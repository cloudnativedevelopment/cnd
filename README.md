# Okteto: A Tool for Cloud Native Developers

[![GitHub release](http://img.shields.io/github/release/okteto/okteto.svg?style=flat-square)][release]
[![CircleCI](https://circleci.com/gh/okteto/okteto.svg?style=svg)](https://circleci.com/gh/okteto/okteto)
[![Apache License 2.0](https://img.shields.io/github/license/okteto/okteto.svg?style=flat-square)][license]

[release]: https://github.com/okteto/okteto/releases
[license]: https://github.com/okteto/okteto/blob/master/LICENSE
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/3055/badge)](https://bestpractices.coreinfrastructure.org/projects/3055)

## Overview

Kubernetes has made it very easy to deploy applications to the cloud at a higher scale than ever, but the development practices have not evolved at the same speed as application deployment patterns.

Today, most developers try to either run parts of the infrastructure locally, or just test these integrations directly in the cluster via CI jobs or the "docker build, docker push, kubectl apply" cycle. It works, but this workflow is painful and incredibly slow.

Okteto makes this cycle a lot faster by letting you launch your development environment directly in your Kubernetes cluster.

## Features

### Development environments on demand 
Your development environment is defined in a [simple yaml manifest](https://okteto.com/docs/reference/manifest).
- Run `okteto init` to inspect your project and generate your own config file 
- Run `okteto up` to launch your development environment in seconds. 

Check in your `okteto.yml` into your repo and make collaboration easier than ever. `git clone`, `okteto up` and you're ready to go.

### Developer Mode 

You can swap your development environment with an existing Kubernetes deployment, and develop directly in a production-like environment. This helps eliminate integration issues since you're developing as if your were in production.

Okteto supports applications with one or with multiple services.

### Kubernetes Live Development

Okteto detects your code changes, synchronizes your code to your development environment. Integrate this with existing compilers and hot reloaders to automatically see your changes in seconds. No commit, build, push or deploy required.

Okteto is powered by [Syncthing](https://github.com/syncthing/syncthing). It will detect code changes instantly and automatically synchronize your file changes. Files are synchronized both ways. If you edit a file directly in your remote development environment, the changes will be reflected locally as well. Great for keeping your `package-lock.json` or `requirements.txt` up to date.

### Keep Your Own Tools
No need to change IDEs, tasks or deployment scripts. Okteto easily integrates and augments with your existing tools.

Okteto is compatible with any Kubernetes cluster, local or remote. If you can run `kubectl apply` you can use Okteto. Our community uses Okteto in all major Kubernetes distros, from Minikube and k3s all the way to GKE, Digital Ocean, AKS, EKS and Civio.

## Learn More
- [How does Okteto work?](docs/how-does-it-work.md)
- Get started following our [installation guides](docs/installation.md).
- Check the Okteto [CLI reference](https://okteto.com/docs/reference/cli) and the [okteto.yml reference](https://okteto.com/docs/reference/manifest).
- Explore our [samples] to learn more about the power of Okteto (https://github.com/okteto/samples).

## Stay in Touch
Got questions? Have feedback? Join [the conversation in Slack](https://kubernetes.slack.com/messages/CM1QMQGS0/)! If you don't already have a Kubernetes slack account, [sign up here](http://slack.k8s.io/). 

Follow [@OktetoHQ](https://twitter.com/oktetohq) on Twitter for important announcements.

Or get in touch with the maintainers:

- [Pablo Chico de Guzman](https://twitter.com/pchico83)
- [Ramiro Berrelleza](https://twitter.com/rberrelleza)
- [Ramon Lamana](https://twitter.com/monchocromo)

## About Okteto
[Okteto](https://okteto.com) is licensed under the Apache 2.0 License.

This project adheres to the Contributor Covenant [code of conduct](code-of-conduct.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to hello@okteto.com.
