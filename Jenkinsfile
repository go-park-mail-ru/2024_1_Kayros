#!/usr/bin/env groovy

pipeline {
  agent any

  stages {
      def microservices = ['comments', 'rest', 'session', 'user', 'auth', 'gateway']
      for (int i = 0; i < microservices.length; i++) {
          stage("Build Microservice: ${microservices[i]}") {
              steps {
                script {
                  sh "cp /home/ubuntu/config.yaml ./config/"
                  sh "sudo docker build -t resto-${microservices[i]}-service:latest -f ./integration/microservices/${microservices[i]}/Dockerfile ."
                }
              }
          }
          
          stage("Push Microservice: ${microservices[i]}") {
              steps {
                script {
                  def localImage = "resto-${microservices[i]}-service:latest"
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
}
