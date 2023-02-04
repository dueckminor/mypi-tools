import os
import yaml
from typing import Union

repo_dir = os.path.normpath(os.path.join(__file__,"../../../../.."))

def mkdir(dirname: str):
    if not os.path.exists(dirname):
        os.makedirs(dirname)

class Config:
    def __init__(self):
        self.mypi_root = os.path.join(repo_dir,".mypi")

    def filename(self, filename:str) -> str:
        return os.path.join(self.mypi_root,filename)

    def write(self, filename:str, content:Union[str,dict]):
        filename = self.filename(filename)
        mkdir(os.path.dirname(filename))
        with open(filename,'w') as stream:
            if isinstance(content,dict):
                yaml.dump(data=content,stream=stream)
            else:
                stream.write(content)
