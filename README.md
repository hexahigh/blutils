# Blutils

[![License](https://img.shields.io/github/license/hexahigh/blutils)](https://github.com/hexahigh/blutils/blob/main/LICENSE)
![Go report card](https://goreportcard.com/badge/github.com/hexahigh/blutils)<br>

A utility program im working on, it is the successor of my last utility program, [boofutils](https://github.com/hexahigh/boofutils).
It contains a few useful utilities and a few things that are just there for fun.

## Download

If you are using a debian based distro you can install blutils using my apt repository.
Simply run these commands:

```bash
echo "deb [signed-by=/usr/share/keyrings/boofdev.apt.pub] https://apt.080609.xyz stable main" | sudo tee -a /etc/apt/sources.list.d/boofdev.list && sudo wget -q -O /usr/share/keyrings/boofdev.apt.pub https://apt.080609.xyz/pgp-key.public
sudo apt update
sudo apt install -y blutils
```