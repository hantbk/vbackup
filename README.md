# vBackup

<p align="center"><a href="https://github.com/hantbk/vbackup" target="_blank"><img src="./web/dashboard/src/assets/logo/vbackup-bar.png" alt="vBackup" width="500" /></a></p>

<p align="center"><b>A Simple, Open-Source, and Modern Web UI for Restic</b></p>

[![Go Report Card](https://goreportcard.com/badge/github.com/hantbk/vbackup)](https://goreportcard.com/report/github.com/hantbk/vbackup)
[![GoDoc](https://godoc.org/github.com/hantbk/vbackup?status.svg)](https://godoc.org/github.com/hantbk/vbackup)
[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go version](https://img.shields.io/github/go-mod/go-version/hantbk/vbackup.svg)](https://github.com/hantbk/vbackup)
<!-- [![GitHub release](https://img.shields.io/github/release/hantbk/vbackup.svg)](https://github.com/hantbk/vbackup/releases/)
[![CI Status](https://github.com/hantbk/vbackup/actions/workflows/ci.yml/badge.svg)](https://github.com/hantbk/vbackup/actions/workflows/ci.yml)
[![Release Status](https://github.com/hantbk/vbackup/actions/workflows/release.yml/badge.svg)](https://github.com/hantbk/vbackup/actions/workflows/release.yml) -->

vBackup is a file backup system built on Restic. It not only has the powerful capabilities of Restic, but also adds rich functions and ease of use to provide users with a more comprehensive data protection experience.

## Features

- **Simple**: The friendly web interface, integrated OTP authentication, and account and password login make the creation, management, and monitoring of backup tasks more intuitive and convenient, making it easy for both novice and experienced users to get started.
- **Efficient**: Based on Restic's incremental backup technology, vBackup only backs up data that has changed since the last backup, effectively saving storage space while maintaining backup speed.
- **Fast**: Written in the go programming language, it fully improves backup efficiency, and the backup data is only limited by the network or hard disk bandwidth.
- **Safety**: Data is fully encrypted during transmission and storage to ensure data confidentiality and integrity, while hash verification is used to ensure data consistency.
- **Compatible**: Supports importing your existing Restic repository.
- **Diversity**: Supports multiple storage backends, including local disk, SFTP, MinIO, AWS, Azure, and more. You can choose the storage backend that best suits your needs.

## Dependencies

- [Restic](https://github.com/restic/restic)
- [Cobra](https://github.com/spf13/cobra)
- [Iris](https://github.com/kataras/iris)
- [Vue](https://github.com/vuejs/vue)
