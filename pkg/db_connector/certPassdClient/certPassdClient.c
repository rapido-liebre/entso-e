#include "certPassdClient.h"
#include <sys/socket.h>
#include <netinet/in.h>
#include <arpa/inet.h>
#include <stdlib.h>
#include <string.h>
#include <sys/types.h>
#include <netdb.h>
#include <stdio.h>
#include <unistd.h>
#include <ctype.h>

#include <errno.h>

char strError[2048];

const unsigned char* BLOWFISH_KEY = (unsigned char*)"l4b!i9xBsn8qKe0G";



bool getPasswordFromServer(const char *name, char *pass, const char *certPassHostSvr, int port)
{
  bool rtn = true;
  if (port <= 0) {
      port = PORT;
  }

  int sockfd, connfd;
  struct sockaddr_in servaddr, cli;
  unsigned char buf[4096];
  unsigned short buflen;
  struct addrinfo hints;
  struct addrinfo *result, *rp;
  in_addr_t addr;
  int s,crv;
  bool connected=false;

  if(name==NULL || pass==NULL){
    sprintf(strError, "name or pass parameters not initialized");
    return false;
  }
  // socket create and varification
  sockfd = socket(AF_INET, SOCK_STREAM, 0);
  if (sockfd == -1) {
      sprintf(strError, "socket()::%s", strerror(errno));
      //printf("socket creation failed...\n");
      return false;
  }
  if(!isdigit(certPassHostSvr[0]) ){
      struct addrinfo hints, *res;
      int errcode;
      void *ptr;
      char addrstr[100];

      memset (&hints, 0, sizeof (hints));
      hints.ai_family = PF_UNSPEC;
      hints.ai_socktype = SOCK_STREAM;
      hints.ai_flags |= AI_CANONNAME;

      errcode = getaddrinfo (certPassHostSvr, NULL, &hints, &res);
      if (errcode != 0){
//          perror ("getaddrinfo");
          sprintf(strError, "getaddrinfo()::%s", gai_strerror(errcode));
          close(sockfd);
          return false;
      } else {
          while (res)
            {
              inet_ntop (res->ai_family, res->ai_addr->sa_data, addrstr, 100);

              switch (res->ai_family)
                {
                case AF_INET:
                  ptr = &((struct sockaddr_in *) res->ai_addr)->sin_addr;
                  break;
                case AF_INET6:
                  ptr = &((struct sockaddr_in6 *) res->ai_addr)->sin6_addr;
                  res = res->ai_next;
                  continue;
                }
              inet_ntop (res->ai_family, ptr, addrstr, 100);
              bzero(&servaddr, sizeof(servaddr));
              addr=inet_addr(addrstr);
              servaddr.sin_family = AF_INET;
              servaddr.sin_addr.s_addr = addr;
              servaddr.sin_port = htons(port);
              if (connect(sockfd, (SA*)&servaddr, sizeof(servaddr)) != 0) {
                  sprintf(strError, "connect(1)::%s", strerror(errno));
                  close(sockfd);
                  return false;
              } else {
                  connected=true;
                  break;
              }
              res = res->ai_next;
            }
      }
  } else {

      bzero(&servaddr, sizeof(servaddr));
      addr=inet_addr(certPassHostSvr);
      servaddr.sin_family = AF_INET;
      servaddr.sin_addr.s_addr = addr;
      servaddr.sin_port = htons(port);
      if (connect(sockfd, (SA*)&servaddr, sizeof(servaddr)) != 0) {
          sprintf(strError, "connect(2)::%s", strerror(errno));
          close(sockfd);
          return false;
      } else {
          connected=true;
      }
  }

  if(connected){
      char tct[1024];
      memset(tct,0,sizeof(tct));
      strcpy(tct,name);
      buflen=strlen(tct);
      memset(buf,0,sizeof(buf));
      memcpy(buf,&buflen,sizeof(unsigned short));
      unsigned char *out=blowfish_encrypt((unsigned char *)tct);
      int rest=buflen%8;
      if(rest >0){
          buflen = buflen-rest+8;
      }
      memcpy(buf+sizeof(unsigned short),out,buflen+1);
      write(sockfd, (void *)buf, buflen+sizeof(unsigned short)+1);
      free(out);
      bzero(buf, sizeof(buf));
      int cread = read(sockfd, buf, sizeof(buf));
      buf[cread]=0;

      unsigned short rlen=cread-sizeof(unsigned short);
      memcpy(&rlen,buf,sizeof(unsigned short));
      rest=rlen%8;
      int modulo8len=rlen;
      if(rest>0){
          modulo8len =  rlen - rest + 8;
      }
      out=blowfish_decrypt(buf+sizeof(unsigned short),modulo8len);

      memcpy(pass,(void *)out,rlen);
      pass[rlen]=0;

      free(out);
      close(sockfd);
      if(rlen>0){
          if(strncmp(pass,name,rlen)==0){
              sprintf(strError, "pass and name are equal");
          }
        return (true);
      }
  }
  close(sockfd);
  sprintf(strError, "unknown error");
  return false;
}

unsigned char* blowfish_decrypt(unsigned char* in, int len){

  int i;
  int SIZE_IN = len;
  unsigned char *out=(unsigned char *)malloc(SIZE_IN+1);

  char ivec[8];
  for(i=0; i<8; i++) ivec[i] = 'i';

  BF_KEY *key = (BF_KEY *)calloc(1, sizeof(BF_KEY));

  /* set up a test key */
  BF_set_key(key, BF_KEY_SIZE, BLOWFISH_KEY );

  BF_cbc_encrypt((unsigned char *)in, (unsigned char *)out, len, key, (unsigned char *)ivec, BF_DECRYPT);

  //printf("Size of out: %d\n", strlen(out));
  //printf("Size of in: %d\n", strlen(in));
  free(key);
  return out;
}

unsigned char* blowfish_encrypt(unsigned char* in){

  int i;
  int SIZE_IN = strlen((char *)in);
  unsigned char *out = (unsigned char *)malloc(SIZE_IN+1);

  char ivec[8];
  for(i=0; i<8; i++) ivec[i] = 'i';

  BF_KEY *key = (BF_KEY *)calloc(1, sizeof(BF_KEY));

  /* set up a test key */
  BF_set_key(key, BF_KEY_SIZE, BLOWFISH_KEY );

  int len=strlen((char *)in);
  int rest=len%8;
  if(rest >0 ){
      len = len-rest+8;
  }

  BF_cbc_encrypt((unsigned char *)in, (unsigned char *)out, len, key, (unsigned char *)ivec, BF_ENCRYPT);

  //printf("Size of out: %d |%s|\n", strlen(out),out);
  //printf("Size of in: %d |%s|\n", strlen(in),in);
  free(key);
  return out;
}
