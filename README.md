# qre

Quick REference in terminal

## Installation from deb

Install with root:
```bash
sudo -E dpkg -i ./deb/qre.deb
export PATH="$PATH:$HOME/.local/bin"
```

Rootless:
```
dpkg -x qre.deb .
mkdir -p ~/.local/bin
mv ./tmp/qre/qres ~/.qre
mv ./tmp/qre/bin/* ~/.local/bin
export PATH="$PATH:$HOME/.local/bin"
rmdir tmp
```

References will be available by
```bash
qre
```

## Building from source

```
make
make install
```
