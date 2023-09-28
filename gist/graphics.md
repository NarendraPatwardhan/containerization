Enabling graphical capabilities on linux requires additional setup such as

1. Sharing the DISPLAY environmental variable
2. Mounting the /tmp/.X11-unix/ directory (which contains the socket) in the
   same location within the container
3. Mounting the Xauthority data (found using Xauthority environmental variable)
   at /root/.Xauthority within the container

Flux takes care of the above settings by default. This enables running many GUI
based applications such as xclock, firefox and more. It also shares all GPUs
from host by default, thus enabling hardware acceleration.

However, running Vulkan based applications such as Rerun or WGPU require
additional setup.

First, based on the
[Github Issue](https://github.com/NVIDIA/nvidia-docker/issues/1480#issuecomment-964285018),
it needs explicitly setting the GPU capabilities. This has become the new
default.

Secondly it needs instllation of the `libvulkan-dev` package and additional
volume mount /usr/share/vulkan/icd.d/nvidia_icd.json in the same location within
container. The user can perform these steps manually as follows:

While spinning up the container add the extra argument

```bash
./flux up ${CONTAINER_NAME} -a "-v /usr/share/vulkan/icd.d/nvidia_icd.json:/usr/share/vulkan/icd.d/nvidia_icd.json"
```

Then once connected and within the container

```bash
sudo apt update && sudo apt install libvulkan-dev
```
