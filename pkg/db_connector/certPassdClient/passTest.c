#include "certPassdClient.h"
#include <stdio.h>
int main( int argc, char** argv )
{
    char pass[MAX];
    char certName[4096]="entso-e:ssir";
    char passSrvAddr[4096]="127.0.0.1";
    int port = 0;

    if(getPasswordFromServer(certName, pass, passSrvAddr, PORT)){
        printf("Sukces hasło dla %s to [%s]\n", certName, pass);
    } else{
        printf("Błąd odczytu hasła dla %s\n", certName);
    }

    return 0;
}
