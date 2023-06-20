#ifndef PASSDCLI_H
#define PASSDCLI_H

#include <sys/socket.h>
#include <openssl/blowfish.h>
#define MAX 80
#define PORT 43099
#define SA struct sockaddr

#define BF_KEY_SIZE 16

#define bool int
#define false 0
#define true 1

unsigned char* blowfish_decrypt(unsigned char* in, int len);
unsigned char* blowfish_encrypt(unsigned char* in);
bool getPasswordFromServer(const char *name, char *pass, const char *certPassHostSvr, int port);


#endif
