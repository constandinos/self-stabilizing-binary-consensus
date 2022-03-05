#!/bin/bash

kill $(ps -u const | grep self-stabiliz | awk '{print $1}')

