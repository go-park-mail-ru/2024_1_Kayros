#!/usr/bin/env groovy

def microservices = 'comment, rest, session, user, auth, gateway'.split(', ')

pipeline {
  agent any 

  tools {go "recent go"}

  stages {
    stage('Initialize') {
    steps {
      script {
      for (int i = 0; i < microservices.length; i++) {
          stage("Build Microservice: ${microservices[i]}") {
                script {
                  sh "sudo cp /home/kayros/backend/config/config.yaml ./config/"
                  sh "sudo docker build -t resto-${microservices[i]}-service:latest -f /var/lib/jenkins/workspace/resto-backend/integration/microservices/${microservices[i]}/Dockerfile ."
                }
          }
      }
      
      stage('Test') {
            script {
                sh 'go test ./... -coverprofile=cover.out'
            }
        }

      stage('Code Analysis') {
          script {
              sh 'sudo curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | bash -s -- -b $GOPATH/bin v1.58.2'
              sh 'golangci-lint run'
          }
      }

      // for (int i = 0; i < microservices.length; i++) {
      //   stage("Push Microservice: ${microservices[i]}") {
      //           script {
      //             def localImage = "resto-${microservices[i]}-service:latest"
      //             def repositoryName = "kayrosteam/${localImage}"
      //             sh "docker tag ${localImage} ${repositoryName} "
      //             docker.withRegistry("", "dockerhub-credentials") {
      //               def image = docker.image("${repositoryName}");
      //               image.push()
      //             }
      //           }
      //     }
      // } 

      }
    }
    }
  }
}
