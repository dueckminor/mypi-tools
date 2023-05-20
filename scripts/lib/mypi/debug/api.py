import os
import yaml
from requests import Session

class API:
    the_api:"API" = None

    @classmethod
    def from_env(cls) -> "API":
        if cls.the_api is not None:
            return cls.the_api
        url=os.environ.get('MYPI_DEBUG_URL')
        secret=os.environ.get('MYPI_DEBUG_SECRET')
        a = API(url=url)
        a.session.get(f'{url}/api/info?local_secret={secret}')
        cls.the_api = a
        return a

    def __init__(self, url:str):
        self.url = url
        self.session = Session()

    def download_file(self,filename:str):
        resp = self.session.get(f'{self.url}/api/fs/opt/mypi/{filename}')
        return resp.text

    def download_yml(self,filename:str) -> dict:
        resp = self.session.get(f'{self.url}/api/fs/opt/mypi/{filename}')
        return yaml.safe_load(resp.text)

    def get_component_info(self, svc:str, comp:str) -> dict:
        resp = self.session.get(f'{self.url}/api/services/{svc}/components/{comp}')
        return yaml.safe_load(resp.text)

    def new_component_port(self, svc:str, comp:str) -> int:
        resp = self.session.patch(
            f'{self.url}/api/services/{svc}/components/{comp}',
            json={'port':0})
        data = yaml.safe_load(resp.text)
        return int(data.get('port'))
    
    def set_component_state(self, svc:str, comp:str, state:str) -> dict:
        resp = self.session.patch(
            f'{self.url}/api/services/{svc}/components/{comp}',
            json={'state':state})
        data = yaml.safe_load(resp.text)
        return data
    
    def set_component_dist(self, svc:str, comp:str, dist:str) -> dict:
        resp = self.session.patch(
            f'{self.url}/api/services/{svc}/components/{comp}',
            json={'dist':dist})
        data = yaml.safe_load(resp.text)
        return data
    