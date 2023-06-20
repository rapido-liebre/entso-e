TARGET = passTest

#include( ../../../../core/core.app )


#CONFIG     += qt
#QT += core
#QT -= gui

#DEFINES +=S_AF_INET_IPC_SOCKET
#HEADERS = $$files( "*.h" )
#SOURCES = passTest.cpp certPassdClient.cpp
#INCLUDEPATH += ../
#LIBS += -lssl -lcrypto


TEMPLATE = app
CONFIG += console
CONFIG -= app_bundle
CONFIG -= qt
HEADERS = certPassdClient.h
SOURCES += \
        passTest.c \
        certPassdClient.c
INCLUDEPATH += ../
LIBS += -lssl -lcrypto

