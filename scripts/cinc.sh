#!/bin/bash

# Derived from https://github.com/moby/moby/blob/master/hack/dind
# Referenced from https://github.com/containerd/containerd/issues/6659

# Allow cgroup nesting
function allow_cgroup_nesting() {
    if [ -f /sys/fs/cgroup/cgroup.controllers ]; then
	# move the processes from the root group to the /init group,
	# otherwise writing subtree_control fails with EBUSY.
	# An error during moving non-existent process (i.e., "cat") is ignored.
	mkdir -p /sys/fs/cgroup/init
	xargs -rn1 < /sys/fs/cgroup/cgroup.procs > /sys/fs/cgroup/init/cgroup.procs || :
	# enable controllers
	sed -e 's/ / +/g' -e 's/^/+/' < /sys/fs/cgroup/cgroup.controllers \
		> /sys/fs/cgroup/cgroup.subtree_control
fi
}

# Cache the invocation to only run once through presence of a log file
function setup() {
    # Make a directory /logs if it doesn't exist
    mkdir -p /logs
    # Check if cinc.log exists within /logs
    if [ ! -f /logs/cinc.log ]; then
        # If it doesn't exist, create it
        touch /logs/cinc.log
        allow_cgroup_nesting
    fi
}

setup

if [[ $# -eq 0 ]]; then
  exec "/bin/bash"
else
  exec "$@"
fi