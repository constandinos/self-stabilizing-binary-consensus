#!/bin/bash

kill $(ps -u $USER | grep main | awk '{print $1}')

