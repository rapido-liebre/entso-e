#!/bin/bash

curl -X POST http://127.0.0.1:3055/api/save_pzrr_publish -H "Content-Type: application/json" -d '{"data":{"creator":"Ernest Malinowski","start":"2022-01","end":"2022-12"},"forecastedCapacityUp":[{"position":1,"quantity":1500},{"position":2,"quantity":1500},{"position":3,"quantity":1500},{"position":4,"quantity":1500}],"forecastedCapacityDown":[{"position":1,"quantity":0},{"position":2,"quantity":0},{"position":3,"quantity":0},{"position":4,"quantity":0}]}'