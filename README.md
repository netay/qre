# qre

Quick REference in terminal

Usage:
```bash
qre
```

## Installation from deb

Install with root:
```bash
sudo -E dpkg -i ./deb/qre.deb
export PATH="$PATH:$HOME/.local/bin"
```

Rootless:
```bash
dpkg -x qre.deb .
mkdir -p ~/.local/bin
mv ./tmp/qre/qres ~/.qre
mv ./tmp/qre/bin/* ~/.local/bin
export PATH="$PATH:$HOME/.local/bin"
rmdir tmp
```

## Building from source

```bash
make
make install
```

