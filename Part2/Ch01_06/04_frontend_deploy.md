## frontend 서비스 구성

![alt text](image-1.png)

1. nginx 설정에 환경 변수 사용을 위해 template 사용
+ nginx/config/nginx.conf 내용 수정
```
events {
  worker_connections  4096;  ## Default: 1024
}
http {
    charset utf-8;
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    include /etc/nginx/conf.d/*.conf;
    include /etc/nginx/sites-enabled/*;
}
```

2. nginx/templates/default.conf.template 파일 생성
 + 디렉토리 생성이 필요한 경우 터미널 창에서 mkdir nginx/templates 수행
```
server {
    listen      80;
    server_name localhost; # IP 또는 FQDN으로 변경. 실습에서는 AWS 콘솔에서 IP 확인
    charset     utf-8;

    # max upload size
    client_max_body_size 75M;   # adjust to taste

    location /static/ {
        alias /data/static/; # Django 프로젝트의 static 파일 위치. 실습에서는 Dockerfile에서 지정.
    }

    location / {
        # 실습에서는 host의 private ip 사용 ($ ip address show eth0)
        #proxy_pass              http://app:8000;  
        proxy_pass              http://poll-backend.${COPILOT_SERVICE_DISCOVERY_ENDPOINT}:8000;  
        proxy_set_header        Host $host;
    }
}
```

3. 프로젝트 루트 위치에 Dockerfile.frontend 파일 생성
```
FROM nginx
COPY nginx/templates /etc/nginx/templates/
COPY nginx/config/nginx.conf /etc/nginx/nginx.conf
COPY static /data/static

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

4. Poll frontend 서비스 생성
```
copilot init

# Workload type: Load Balanced Web Service
# Service name: poll-frontend
# Dockerfile: ./Dockerfile.frontend
```
coplit init -> Load Balanced Web Service -> name : poll-frontend 지정

1. ELB 헬스체크를 지원하기 위해 헬스체크용 미들웨어 구성
- polls/middleware.py 파일 생성
```
from django.http import HttpResponse


class HealthCheckMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        if request.path == "/health":
            return HttpResponse("ok")
        response = self.get_response(request)
        return response
```

- pollme/settings.py 수정
```
MIDDLEWARE = [
    'polls.middleware.HealthCheckMiddleware',
    'django.middleware.security.SecurityMiddleware',
    'django.contrib.sessions.middleware.SessionMiddleware',
...
```

6. poll-frontend 서비스에 health check 경로 설정
 + copilot/poll-frontend/manifest.yml 내용 수정
```
http:
  # Requests to this path will be forwarded to your service.
  # To match all requests you can use the "/" path.
  path: '/'
  # You can specify a custom health check path. The default is "/".
  healthcheck: '/health'
```

7. Poll frontend 서비스 배포
copilot deploy


8. AWS Web Console에서 Poll frontend 서비스 배포 상태 확인 
+ 400 오류 확인 --> 트러블 슈팅

9. 배포된 Poll 앱 확인 
+ Invalid HTTP_HOST header 오류 확인 --> 트러블 슈팅

10. pollme/settings.py 수정
![alt text](image-2.png)
```
# SECURITY WARNING: don't run with debug turned on in production!
DEBUG = False

ALLOWED_HOSTS = ['.elb.amazonaws.com'] # host allow 필요 (Django 보안)
```

11. Poll backend 서비스 빌드/배포
copilot deploy

12. Test
+ 배포된 Poll 앱 확인 
+ 테스트: 회원 가입 및 로그인 
