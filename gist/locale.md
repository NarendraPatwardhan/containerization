Setting up a locale inside a docker container is important to ensure that the
applications contained within will behave consistently and will be able to
handle text input and output correctly. This is particularly important for
development images which form the bulk of this repository.

In the past, the locale setup was done using the following snippet, immediately
after setting up the parent image and headless installation commands.

```dockerfile
RUN apt update && \
    $APT_INSTALL locales && \
    locale-gen en_US.UTF-8 && \
    echo LANG=en_US.UTF-8 > /etc/default/locale

ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8
```

This installs the locale package, generates configuration for English UTF-8, and
then sets it as the default locale. It also sets the environmental variables
which might be queried by other programs. It also sets the environmental
variables which might be queried by other programs. The LC_ALL is the variable
that has the highest priority for checking but does not need to be set as
programs fall back to individual LC_XX variables that are better suited. The
LANGUAGE variable is used to display the man pages in the order of preference.

[This video](https://www.youtube.com/watch?v=kL0q-7alfQA) by
[Anthony Sottile](https://twitter.com/codewithanthony), explains that simply
setting up the LANG variable to the C.UTF-8 value can be a better solution.
Moving forward all images shall set the locale up the following way for
simplicity, however, if a particular application errors out, feel free to
regress or use the above snippet to manually make a change in the development
environment.

```dockerfile
ENV LANG=C.UTF-8
```

The change has been tested for behavioral equivalence in VS-Code and Neovim
using the Python `emoji` package.
