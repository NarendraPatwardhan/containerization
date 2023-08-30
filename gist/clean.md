Dockerfiles store each instruction as a layer. Thus if the dockerfile contains
two instructions as follows:

```dockerfile
# Create a file containing 100 MB of random data
RUN dd if=/dev/urandom of=random_data.bin bs=1M count=100
# Remove the data
RUN rm -rf random_data.bin
```

The end image will retain the 100 MB file in the first layer. The second layer
simply creates a white-out file to indicate the file has been deleted. This
leads to inflated image sizes.

To avoid this, both instructions can be combined in one RUN command

```dockerfile
RUN dd if=/dev/urandom of=random_data.bin bs=1M count=100 && \
    rm -rf random_data.bin
```

This avoids storing the file in the layer and also removes the extra layer.

This practice has been followed since the start of this repository for any
downloaded or generated files, however, one of the "best" practices builds on
this and tells you to similarly clean apt lists and clean cache after every apt
install. No perceivable difference beyond a couple of KBs per instruction was
found after incorporating the APT_CLEAN. Thus we retain our original formulation
for the sake of avoiding complexity.
