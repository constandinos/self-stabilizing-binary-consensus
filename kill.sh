#!/bin/bash

kill $(ps -u $USER | grep self-stabiliz | awk '{print $1}')

