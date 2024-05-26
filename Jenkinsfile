#!/usr/bin/env groovy

pipeline {
  agent any
  stages {
     stage("Build Microservice: Auth") {
        steps {
           script {
              sh "cp /home/ubuntu/config.yaml ./config/"
              sh "sudo docker build -t resto-auth-service:latest -f ./integration/microservices/auth/Dockerfile ."
           }
        }
      }

      stage("Push to Dockerhub") {
          steps {
            script {
                def localImage = "resto-auth-service:latest"
                def repositoryName = "kayrosteam/${localImage}"
                sh "docker tag ${localImage} ${repositoryName} "
                docker.withRegistry("", "dockerhub-credentials") {
                  def image = docker.image("${repositoryName}");
                  image.push()
                }
            }
          }
        }
  }
}
