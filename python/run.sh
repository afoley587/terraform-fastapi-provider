#!/bin/bash

export DES_KEY='hello123'

poetry run uvicorn app:app --reload
