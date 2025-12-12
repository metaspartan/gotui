# SSH Dashboard Example

This example demonstrates the **SSH Dashboard** widget/feature.

## 🚀 Run

```bash
$ ssh-keygen -t ed25519 -f hostkey -N "" # Generate sample host key
$ go run _examples/stacked_barchart/main.go
```

In a separate window:

```bash
$ ssh 0.0.0.0 -p 2222
```

## 📝 Code

See [main.go](main.go) for the implementation.
