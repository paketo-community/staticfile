# Staticfile Cloud Native Buildpack

The buildpack automates the process of creating configuration files for nginx. 

## Integration

```toml
[[requires]]

  # The name of the Staticfile dependency is "staticfile". This value is considered
  # part of the public API for the buildpack and will not change without a plan
  # for deprecation.
  name = "staticfile"
```

## Usage

To package this buildpack for consumption:

```
$ ./scripts/package.sh
```

This builds the buildpack's Go source using `GOOS=linux` by default. You can supply another value as the first argument to `package.sh`.


## `buildpack.yml` Configurations

```yaml
staticfile:
  nginx:
    root:
    host_dot_files:
    location_include:
    directory:
    ssi:
    pushstate:
    http_strict_transport_security:
    http_strict_transport_security_include_subdomains:
    http_strict_transport_security_preload:
    force_https:
    basic_auth:
    status_codes:
```

As of now this buildpack only integrates with `nginx`, so to get a default nginx config simply leave the map under `nginx` empty.

```yaml
staticfile:
  nginx: {}
```

If you do not specify `nginx`, the `detect` binary will return an error.
