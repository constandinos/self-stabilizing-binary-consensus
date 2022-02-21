#!/bin/bash

kill $(ps | grep self-stabiliz | awk '{print $1}')

