#!/usr/bin/env groovy

pipeline {
  agent any
  stages {
     stage("Build Microservice: Auth") {
        steps {
           script {
              sh "cp /home/ubuntu/config.yaml ./config/"
              //docker.build("resto-auth-service:latest", "-f ./integration/microservices/auth .")
           }
        }
      }
  }
}
