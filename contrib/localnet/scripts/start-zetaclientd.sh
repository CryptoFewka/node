#!/bin/bash
#DISCOVERED_HOSTNAME=$(hostname)
DISCOVERED_HOSTNAME=$(nslookup $(hostname -i) | grep '=' | awk -F'= ' '{split($2, a, "."); print a[1]}')
DISCOVERED_NETWORK=$(echo $DISCOVERED_HOSTNAME |  awk -F'-' '{split($1, a, "-"); print a[1]}')

cp  /root/preparams/PreParams_$DISCOVERED_HOSTNAME.json /root/preParams.json
num=$(echo $DISCOVERED_HOSTNAME | tr -dc '0-9')
node="zetacore_$num"
mv  /root/zetacored/zetacored_$node /root/.zetacored

mv /root/tss/$DISCOVERED_HOSTNAME /root/.tss


if [ $DISCOVERED_HOSTNAME == "$DISCOVERED_NETWORK-zetaclient-1" ]
then
    rm ~/.tss/address_book.seed
    zetaclientd init \
      --pre-params ~/preParams.json  --zetacore-url $DISCOVERED_NETWORK-zetacore-1 \
      --chain-id athens_101-1 --operator zeta1z46tdw75jvh4h39y3vu758ctv34rw5z9kmyhgz --log-level 1 --hotkey=val_grantee_observer
    zetaclientd start
else
  num=$(echo $DISCOVERED_HOSTNAME | tr -dc '0-9')
  node="zetacore-$num"
  SEED=$(curl --retry 10 --retry-delay 5 --retry-connrefused  -s $DISCOVERED_NETWORK-zetaclient-1:8123/p2p)
  ZETACLIENT_IP=$(dig +short $DISCOVERED_NETWORK-zetaclient-1 | awk '{ print; exit }')
  zetaclientd init \
    --peer /ip4/"$ZETACLIENT_IP"/tcp/6668/p2p/$SEED \
    --pre-params ~/preParams.json --zetacore-url $DISCOVERED_NETWORK-$node \
    --chain-id athens_101-1 --operator zeta1lz2fqwzjnk6qy48fgj753h48444fxtt7hekp52 --log-level 0 --hotkey=val_grantee_observer
  zetaclientd start
fi
