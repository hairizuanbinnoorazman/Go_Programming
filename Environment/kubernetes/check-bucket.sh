#!/bin/sh
pod=$(kubectl get pods | grep $1 | awk '{ print $1 }')
state=$(kubectl get pods | grep $1 | awk '{ print $3 }')

if [ -z "$pod" ];
then 
  echo 'No such pod'
  exit 1
fi

if [ $state = "Completed" ];
then
  echo Bucket $1 created successfully
  exit 0
else
  echo Error
  kubectl logs $pod
  exit 1
fi