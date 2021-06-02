[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/go-config.svg)](https://pkg.go.dev/github.com/tommzn/aws-dynamodb)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/aws-dynamodb)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/aws-dynamodb)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/aws-dynamodb)](https://goreportcard.com/report/github.com/tommzn/aws-dynamodb)
[![Actions Status](https://github.com/tommzn/aws-dynamodb/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/aws-dynamodb/actions)

# DynamoDb Wrapper
This package provides a wrapper to dynamodb to run CRUD actions for items. It expects a table with a composed primary key of Id (Hash) and ObjectType (Range).
Sub package testing will support you creating a suitable table for tests.
