#!/usr/bin/python
import os
import shlex
import subprocess

LOCAL_PATH = os.path.abspath(os.path.dirname(__file__))
DOC_PATH = os.path.join(LOCAL_PATH, '..', 'doc')

# Check if the "doc" folder is present where we can fetch all of the godot
# documentation to include it as Go documentation as part of class generation.
if not os.getenv('NODOC'):
    if os.path.exists(DOC_PATH):
        print('Godot documentation found. Pulling the latest changes...')
        os.chdir(DOC_PATH)
        subprocess.run(shlex.split('git pull origin master'))
        os.chdir(LOCAL_PATH)
    else:
        print('Godot documentation not found. Pulling the latest changes...')
        os.makedirs(DOC_PATH)
        os.chdir(DOC_PATH)
        cmds =  [
            'git init',
            'git remote add -f origin https://github.com/godotengine/godot.git',  # noqa
            'git config core.sparseCheckout true',
            'echo "doc/classes" >> {}'.format(
                os.path.join(DOC_PATH, '.git', 'info', 'sparse-checkout')
            ),
            'git pull origin master'
        ]

        for cmd in cmds:
            subprocess.run(shlex.split(cmd))

        os.chdir(LOCAL_PATH)

# clean the gdnative directory
print('Cleaning previous generation...')
subprocess.run(shlex.split('rm gdnative/*.gen.*'), capture_output=True)

env = os.environ.copy()
env.update({'API_PATH': os.path.join(LOCAL_PATH, '..')})
cmd = 'go run -v {}'.format(os.path.join(LOCAL_PATH, 'generate', 'main.go'))
subprocess.run(shlex.split(cmd), env=env)
