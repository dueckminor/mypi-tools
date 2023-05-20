import os
import pyjson5
import json

from typing import Optional,List
from .config import repo_dir

class LaunchJson:
    def __init__(self, filename:Optional[str]=None):
        if filename:
            self.filename = filename
        else:
            self.filename = os.path.join(repo_dir,'.vscode','launch.json')
        self.data:dict
        self.configurations:List[dict]
        self.load()

    def load(self):
        if os.path.exists(self.filename):
            with open(self.filename,'r') as f:
                self.data = pyjson5.decode_io(f)
        else:
            self.data = dict()
        if not self.data.get('version'):
            self.data['version']='0.2.0'
        configurations = self.data.get('configurations')
        if not configurations:
            configurations = []
            self.data['configurations']=configurations
        self.configurations = configurations

    def save(self):
        with open(self.filename,'w') as f:
            json.dump(self.data,f,indent=2)

    def set_configuration(self, configuration:dict):
        name = configuration['name']
        for i in range(len(self.configurations)):
            if self.configurations[i].get('name')==name:
                self.configurations[i] = configuration
                return
        self.configurations.append(configuration)


