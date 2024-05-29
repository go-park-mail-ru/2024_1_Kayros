#!/usr/bin/env groovy

def microservices = 'comment, rest, session, user, auth, gateway'.split(', ')

pipeline {
    agent any 
    tools {go "recent go"}

    stages {
      stage('Initialize') {
          steps {
              script {
                 sh "sudo cp /home/kayros/backend/config/config.yaml ./config/"
                 for (int i = 0; i < microservices.length; i++) {
                     def currentChange
                     def externalChange

                     for (change in currentBuild.changeSets) {
                        for (entry in change.getItems()) {
                            for (file in entry.getAffectedFiles()) {
                               if (file.getPath() =~ "${microservices[i]}") {
                                   currentChange = 'true'
                               }
                             }
                        }
                      }

                      for (service in microservices) {
                          for (change in currentBuild.changeSets) {
                              for (entry in change.getItems()) {
                                  for (file in entry.getAffectedFiles()) {
                                      if (file.getPath().contains(service)) {
                                          externalChange = 'true'
                                      }
                                   }
                               }
                           }
                       }

                       if (currentChange || !externalChange) {
                          stage("Build Microservice: ${microservices[i]}") {        
                              script {
                                 // sh "sudo cp /home/kayros/backend/config/config.yaml ./config/"
                                 // sh "sudo docker build -t resto-${microservices[i]}-service:latest -f ./integration/microservices/${microservices[i]}/Dockerfile ."
                                  def service = microservices[i]
                                  if (microservices[i] != 'gateway') {
                                    service = "${microservices[i]}-grpc"
                                  }

                                  sh "sudo docker-compose -f ./integration/prod-compose/docker-compose.yaml up -d --no-deps --build ${service}"
                              }
                          }
                       }
                 }

                // if (params.FULL_BUILD && params.FULL_BUILD == true) {
                //    stage("Full Project Build") {        
                //         script {
                //             sh "sudo docker-compose -f ./integration/prod-compose/docker-compose.yaml up -d --no-deps --build"  
                //         }
                //    }
                // }
      
                 stage('Test') {
                    script {
                       // sh 'go test ./... -coverprofile=cover.out'
                    }
                 }

                 stage('Code Analysis') {
                    script {
                         sh 'make easyjs'
                           //     sh '/home/kayros/go/bin/golangci-lint run'
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
