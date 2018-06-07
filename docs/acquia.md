Acquia
=======

## Overview

Acquia provides a Cloud API for developers to download database backups.

https://cloudapi.acquia.com

The following command downloads and extracts the compressed backup file.

## BE CAFEFUL

**Problem**

These backup still need to be sanitized!

**Solution**

* Import the database into a temp MySQL database and export using `mtk db dump` command
* Sanitize an Acquia nonprod environment so the backups are sanitized ONLY for that environment and can be used by developers

## Requirements

* Acquia Account

## Usage

```bash
mtk acquia dump --username=myusername \
                --password=mypassword \
                --site=prod:mysite \
                --environment=test \
                --database=testdb \
                --file=/tmp/db.sql
```