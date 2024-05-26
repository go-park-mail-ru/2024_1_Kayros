#!/usr/bin/env groovy

pipeline {
  agent any
  stages {
     stage("Build Microservice: Auth") {
        steps {
           script {
              sh "cp /home/ubuntu/config.yaml ./config/"
              sh "sudo docker build -t resto-auth-service:latest -f ./integration/microservices/auth ."
           }
        }
      }
  }
}
