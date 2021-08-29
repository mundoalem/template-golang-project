<div align="center">

# Template Golang Project

![Go Version](https://img.shields.io/github/go-mod/go-version/mundoalem/template-golang-project)
![Release Version](https://img.shields.io/github/v/release/mundoalem/template-golang-project)
![Pipeline Status](https://github.com/mundoalem/template-golang-project/actions/workflows/pipeline.yml/badge.svg)
[![codecov](https://codecov.io/gh/mundoalem/template-golang-project/branch/main/graph/badge.svg?token=R0HJ0SAOC0)](https://codecov.io/gh/mundoalem/template-golang-project)
![Contributors](https://img.shields.io/github/contributors/mundoalem/template-golang-project)

A DevOps centric template to bootstrap Go projects.

</div>

## Introduction

This project is a template anyone can use in order to bootstrap a project using the Go
programming language. The template has a full feature pipeline following the latest DevOps
practices.

You can find [here](https://github.com/golang-standards/project-layout)
a full explanation of the directory structure.

You can control the project through the use of `[mage](https://magefile.org/)`, it accepts
the following targets:

| Argument | Description                                                       |
| -------- | --------------------------------------------------------------    |
| build    | Runs `go build` saving the final assets in `build` directory      |
| clean    | Removes compiled binaries generated by the build step             |
| lint     | Runs `go fmt` for `cmd`, `internal` and `pkg` directories         |
| lock     | Install dependencies from `go.mod`                                |
| release  | Creates a package and, if in pipeline, also creates a release (1) |
| reset    | Removes installed dependencies                                    |
| scan     | Runs vulnerability scans                                          |
| test     | Runs all tests                                                    |

1. The package will be created for the supported platforms

## License

[GPLv3](https://choosealicense.com/licenses/gpl-3.0/)

## Tech Stack

These are the software baked in this template:

- [Go](https://www.python.org/) 1.16.5
- [Cli](github.com/mitchellh/cli)


## Usage

```bash
$ mage lint
$ mage test
$ mage build
$ mage release
```

## Environment Variables

| Variable      | Description                                                           |
| ------------- | --------------------------------------------------------------------- |
| CODECOV_TOKEN | Token used to calculate and report test coverage from codecov.io      |
| SNYK_TOKEN    | Token used during the security vulnerabilities scan task from snyk.io |

In the pipeline the build script also looks into a few environment variables that are set by
GitHub automatically. These variables are used by the automation in order to either make
decisions or use as input data:

| Variable   | Description                                                           |
| ---------- | --------------------------------------------------------------------- |
| CI         | Used to determine whether we are running locally or inside a pipeline |
| GITHUB_SHA | Used whenever we need to know the current commit hash                 |


## Feedback

If you have any feedback, please open an [issue](https://github.com/mundoalem/template-golang-project/issues).


## Contributing

Please contribute to this repository if any of the following is true:

- You have expertise in community development, communication, or education
- You want open source communities to be more collaborative and inclusive
- You want to help lower the burden to first time contributors

The Prerequisites to contribute:

- Familiarity with [pull requests](https://help.github.com/articles/using-pull-requests) and [issues](https://guides.github.com/features/issues/).
- Knowledge of [Markdown](https://help.github.com/articles/markdown-basics/) for editing .md documents.
- Knowledge of [Go](https://golang.org/) and its ecosystem.

In particular, this community seeks the following types of contributions:

- Ideas: participate in an issue thread or start your own to have your voice heard.
- Resources: submit a pull request to add to RESOURCES.md with links to related content.
- Outline sections: help us ensure that this repository is comprehensive. if there is a topic
  that is overlooked, please add it, even if it is just a stub in the form of a header and
  single sentence. Initially, most things fall into this category.
- Writing: contribute your expertise in an area by helping us expand the included content.
- Copy editing: fix typos, clarify language, and generally improve the quality of the content.
- Formatting: help keep content easy to read with consistent formatting.
- Features: add new features to the project.
- Bugfixes: fix open issues.


## Conduct

We are committed to providing a friendly, safe and welcoming environment for all, regardless of
gender, sexual orientation, disability, ethnicity, religion, income or similar personal
characteristic.

Please be kind and courteous. There's no need to be mean or rude. Respect that people have
differences of opinion and that every design or implementation choice carries a trade-off and
numerous costs. There is seldom a right answer, merely an optimal answer given a set of values
and circumstances.

Please keep unstructured critique to a minimum. If you have solid ideas you want to experiment
with, make a fork and see how it works.

We will exclude you from interaction if you insult, demean or harass anyone. That is not welcome
behavior. We interpret the term "harassment" as including the definition in the
(Citizen Code)[http://citizencodeofconduct.org/] of Conduct; if you have any lack of clarity
about what might be included in that concept, please read their definition. In particular,
we don't tolerate behavior that excludes people in socially marginalized groups.

Whether you're a regular contributor or a newcomer, we care about making this community a safe
place for you and we've got your back.

Likewise any spamming, trolling, flaming, baiting or other attention-stealing behavior is not
welcome.


## Communication

GitHub issues are the primary way for communicating about specific proposed changes to this project.

In both contexts, please follow the conduct guidelines above. Language issues are often contentious and we'd like to keep discussion brief, civil and focused on what we're actually doing, not wandering off into too much imaginary stuff.


## FAQ

### Will there ever be support for other continuous integration platforms?

Right now I have no plans to support other platforms like TravisCI, CircleCI or
Gitlab. Anyway, it should be quite easy for you to port the GitHub Actions to
any platform you like.

The reason for that is that I don't want to have a `.travis.yml`, a
`circleci.yml` and a `.gitlab-ci.yml` all together in the same place when only
one would actually be used. So I want to avoid (for now) cluttering the
template with too many files that might or might not be useful.
