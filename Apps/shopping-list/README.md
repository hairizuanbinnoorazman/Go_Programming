# Shopping List mini application

The following is simple small golang application that is to be deployed on Google Cloud Run and uses datastore as the form of its storage

This is a really small application so various "concurrent" based operations is not supported. We need to reduce any "parallel" actions.

## Operating

We need to generate a couple of keys for "security" purposes. Generate hashkey and blockkey via the following command.

```
echo $RANDOM | md5sum
```