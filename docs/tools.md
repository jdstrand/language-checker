# Tools

## Pre-commit hook

`language-checker` supports being run from a [pre-commit](https://pre-commit.com/) hook,
allowing you to avoid accidentally committing uses of non-inclusive
language.
You have two alternative mechanisms for doing so.
If you have arranged to install `language-checker` on your command search path (as well
as anyone working on your repository), then add this configuration to your
`.pre-commit-config.yaml`:

```yaml
-   repo: https://github.com/jdstrand/language-checker
    rev: ''  # pick a tag to point to
    hooks:
    -   language-checker
```

(Note that in this case the `rev` only controls the version of a wrapper
script that is used, not the version of `language-checker` itself.)

Alternatively, you can tell `pre-commit` to build `language-checker` from source,
although this requires you and anyone working on your repository to have
`go` on your command search path and for it to be at least version 1.18:

```yaml
-   repo: https://github.com/jdstrand/language-checker
    rev: ''  # pick a tag to point to
    hooks:
    -   language-checker-from-source
```

(In this case the `rev` controls the version of `language-checker` itself.)

See the [pre-commit
documentation](https://pre-commit.com/#pre-commit-configyaml---hooks) for
how to customize this further.
