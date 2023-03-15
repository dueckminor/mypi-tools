import subprocess
import socket
import os
import glob
from typing import Optional,List
from .api import API
from .config import repo_dir
import importlib
from requests import get
from threading import Thread
from time import sleep

#import importlib.machinery
#import importlib.util

class Ctrl:
    def __init__(self, service:str, component: str):
        self.service = service
        self.component = component

    def run(self):
        """
        run starts a component and waits until it completes
        """
        pass

    def set_state(self, state:str):
        a = API.from_env()
        a.set_component_state(self.service,self.component,state)

    def get_port(self) -> int:
        a = API.from_env()
        return a.new_component_port(self.service,self.component)

    def wait_for_port(self, component:str) -> int:
        a = API.from_env()
        while True:
            info = a.get_component_info(self.service,component)
            if 'running' == info.get('state'):
                port = int(info.get('port'))
                if port > 0:
                    return port
            sleep(0.5)


    @classmethod
    def from_file(cls, filename:str) -> "Ctrl":
        component = os.path.basename(os.path.dirname(filename))
        service = os.path.basename(os.path.normpath(os.path.join(filename,"../../..")))
        if component == 'web':
            return CtrlWeb(service)
        if component == 'go':
            return CtrlGo(service)

    @classmethod
    def load(cls, service:str, component:str) -> "Ctrl":
        ctrl_file = os.path.join(repo_dir,
            "debug","services",service,"components",component,"ctrl")

        service=service.replace('-','_')
        module_name="mypi.services."+service.replace('-','_')+".components."+component.replace('-','_')+".ctrl"

        loader = importlib.machinery.SourceFileLoader( fullname=module_name, path=ctrl_file )
        spec = importlib.util.spec_from_loader( module_name, loader )
        module = importlib.util.module_from_spec( spec )
        loader.exec_module( module )

        return module.ctrl

class WaitForPort(Thread):
    def __init__(self, ctrl:Ctrl, port:int):
        Thread.__init__(self)
        self.ctrl = ctrl
        self.port = port
        self.stopped = False
        self.start()

    def run(self):
        while not self.stopped:
            try:
                resp = get(f'http://localhost:{self.port}/index.html',timeout=5)
                if resp.status_code >= 200 and resp.status_code < 300:
                    print("RUNNING")
                    self.ctrl.set_state("running")
                    break
            except Exception:
                pass
            sleep(0.5)

class WaitForRawPort(Thread):
    def __init__(self, ctrl:Ctrl, port:int):
        Thread.__init__(self)
        self.ctrl = ctrl
        self.port = port
        self.stopped = False
        self.start()

    def run(self):
        while not self.stopped:
            try:
                with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
                    s.connect(("localhost", self.port))
                    print("RUNNING")
                    self.ctrl.set_state("running")
                    break
            except Exception:
                pass
            sleep(0.5)

class CtrlWeb(Ctrl):
    def __init__(self, service:str):
        Ctrl.__init__(self, service=service, component="web")

    def run(self):
        self.set_state("starting")

        port = self.get_port()

        w = WaitForPort(self,port)

        cwd=f'{repo_dir}/web/{self.service}'
        subprocess.run(args=['npm','install'],cwd=cwd)
        subprocess.run(args=
            [
                os.path.join(cwd,'node_modules/.bin/vue-cli-service'),
                'serve',
                '--host', 'localhost',
                '--port', str(port)
            ],cwd=cwd)

        w.stopped = True
        self.set_state("stopped")


class CtrlGo(Ctrl):
    def __init__(self, service:str, web:bool=True):
        Ctrl.__init__(self, service=service, component="go")
        self.web = web

    def run(self):
        self.set_state("starting")

        if self.web:
            web_port=self.wait_for_port("web")

        cwd=f'{repo_dir}/cmd/{self.service}'

        go_files = glob.glob(f'{cwd}/*.go')

        args=['go','run',]
        args.extend(go_files)
        args.append('--localhost-only')
        if self.web:
            args.append(f'--webpack-debug=http://localhost:{web_port}')

        args.extend(self._get_port_args())

        w = WaitForRawPort(self,self.get_port())
        proc = subprocess.run(args=args,cwd=cwd)
        w.stopped = True
        print(f'RC: {proc.returncode}')

        self.set_state("stopped")

    def _get_port_args(self) -> List[str]:
        return ["--port",str(self.get_port())]
