# entso-e

Compile "C" code into golang:
https://www.sobyte.net/post/2022-06/go-c/

Compile the static library first, generating the libcertPassdClient.a static link library in the certPassdClient directory.

```cd /db_connector/certPassdClient/
gcc -c certPassdClient.c -o certPassdClient.o  -lssl -lcrypto 

ar -crs libcertPassdClient.a certPassdClient.o

certPassdClient.c
certPassdClient.h
certPassdClient.o
libcertPassdClient.a
```

Test for get stored password:
```
gcc passTest.c -o test -L./ -lcertPassdClient -lssl -lcrypto
./test
```

Manage user/password:
```
cd RTDB/dr3/dr3_workspace/dr3x64/core/LINUX64/bin/
certpasswd --fname ~/certPassd.cfg --name ssir --secret ssir
certPassd --memory ~/RTDB/rgx_py_mods/rgx3_pyrtdb/devel_test/SH_RTDB --log ~/certPassd.log --secrets ~/certPassd.cfg
```
