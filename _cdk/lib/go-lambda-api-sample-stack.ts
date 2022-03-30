import { RemovalPolicy, Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
// import * as sqs from 'aws-cdk-lib/aws-sqs';
import { aws_apigateway as apigateway } from "aws-cdk-lib";
import { Code, Function, Runtime } from 'aws-cdk-lib/aws-lambda';
import { StageContext } from './context';
import { AttributeType, Table } from 'aws-cdk-lib/aws-dynamodb';
import { AccountRecovery, OAuthScope, UserPool, UserPoolEmail } from 'aws-cdk-lib/aws-cognito';

export class GoLambdaApiSampleStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const context: StageContext = this.context

    /** User Pool */
    const userPool = new UserPool(this, "test-user-pool", {
      selfSignUpEnabled: true,
      signInCaseSensitive: false,
      autoVerify: { email: true },
      userPoolName: "test-user-pool",
      standardAttributes: {
        familyName: { required: true, mutable: true },
        givenName: { required: true, mutable: true },
        email: { required: true, mutable: true }
      },
      userInvitation: {
        emailSubject: "test",
        // emailBody: "created user {username} {####}"
      },
      email: UserPoolEmail.withSES({
        sesRegion: "ap-northeast-1",
        fromEmail: context.ses.fromEmail
      }),
      accountRecovery: AccountRecovery.EMAIL_ONLY,
      removalPolicy: RemovalPolicy.DESTROY
    })
    userPool.addClient("test-client", {
      userPoolClientName: "test-client",
      oAuth: {
        scopes: [
          OAuthScope.OPENID,
          OAuthScope.OPENID,
          OAuthScope.PROFILE,
        ]
      },
      authFlows: {
        userPassword: true
      }
    })

    /** functions **/
    const listUsers = new Function(this, "listUsers", {
      description: "list users",
      code: Code.fromAsset('../bin/zip/list-users.zip'),
      handler: 'list-users',
      runtime: Runtime.GO_1_X,
    })
    const createUser = new Function(this, "createUser", {
      description: "create user",
      code: Code.fromAsset('../bin/zip/create-user.zip'),
      handler: 'create-user',
      runtime: Runtime.GO_1_X,
    })
    const getUser = new Function(this, "getUser", {
      description: "get user by id",
      code: Code.fromAsset('../bin/zip/get-user.zip'),
      handler: 'get-user',
      runtime: Runtime.GO_1_X,
    })
    const deleteUser = new Function(this, "deleteUser", {
      description: "delete user by id",
      code: Code.fromAsset('../bin/zip/delete-user.zip'),
      handler: 'delete-user',
      runtime: Runtime.GO_1_X,
    })

    /** API **/
    const api = new apigateway.RestApi(this, "listUsersAPI", {
      restApiName: "testAPI",
      endpointTypes: [apigateway.EndpointType.REGIONAL],
      deployOptions: {
        stageName: context.api.stageName
      }
    })
    const users = api.root.addResource("users")
    users.addMethod("GET",
      new apigateway.LambdaIntegration(listUsers, {
        proxy: true,
        passthroughBehavior: apigateway.PassthroughBehavior.WHEN_NO_TEMPLATES,
        integrationResponses: [
          {
            statusCode: "200",
            responseTemplates: { "application/json": "" },
          },
          {
            statusCode: "400",
            responseTemplates: { "application/json": "" },
          },
          {
            statusCode: "500",
            responseTemplates: { "application/json": "" },
          },
        ]
      }),
      {
        methodResponses: [
          {
            statusCode: "200",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "400",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "500",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          }
        ]
      })
    users.addMethod("POST",
      new apigateway.LambdaIntegration(createUser, {
        proxy: true,
        passthroughBehavior: apigateway.PassthroughBehavior.WHEN_NO_TEMPLATES,
        requestTemplates: {
          "application/json": `{ "user": "$input.json('$.user')" }`
        },
        integrationResponses: [
          {
            statusCode: "201",
            responseTemplates: { "application/json": "" },
          },
          {
            statusCode: "400",
            responseTemplates: { "application/json": "" },
          },
          {
            statusCode: "409",
            responseTemplates: { "application/json": "" },
          },
          {
            statusCode: "500",
            responseTemplates: { "application/json": "" },
          }
        ]
      }),
      {
        methodResponses: [
          {
            statusCode: "201",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "400",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "409",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "500",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          }
        ]
      })


    const user = users.addResource("{user_id}")
    user.addMethod("GET",
      new apigateway.LambdaIntegration(getUser, {
        proxy: true,
        passthroughBehavior: apigateway.PassthroughBehavior.WHEN_NO_TEMPLATES,
        requestTemplates: {
          "application/json": `{ "userID": "$input.params('user_id')" }`
        },
        requestParameters: {
          "integration.request.path.userID": "method.request.path.userID"
        },
        integrationResponses: [
          {
            statusCode: "200",
            responseTemplates: { "application/json": "user_id" },
          },
          {
            statusCode: "400",
            responseTemplates: { "application/json": "user_id" },
          },
          {
            statusCode: "404",
            responseTemplates: { "application/json": "user_id" },
          },
          {
            statusCode: "500",
            responseTemplates: { "application/json": "user_id" },
          }
        ]
      }),
      {
        requestParameters: {
          "method.request.path.userID": true
        },
        methodResponses: [
          {
            statusCode: "200",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "400",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "404",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "500",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          }
        ]
      })
    user.addMethod("DELETE",
      new apigateway.LambdaIntegration(deleteUser, {
        proxy: true,
        passthroughBehavior: apigateway.PassthroughBehavior.WHEN_NO_TEMPLATES,
        requestTemplates: {
          "application/json": `{ "userID": "$input.params('user_id')" }`
        },
        requestParameters: {
          "integration.request.path.userID": "method.request.path.userID"
        },
        integrationResponses: [
          {
            statusCode: "204",
            responseTemplates: { "application/json": "user_id" },
          },
          {
            statusCode: "400",
            responseTemplates: { "application/json": "user_id" },
          },
          {
            statusCode: "404",
            responseTemplates: { "application/json": "user_id" },
          },
          {
            statusCode: "500",
            responseTemplates: { "application/json": "user_id" },
          }
        ]
      }),
      {
        requestParameters: {
          "method.request.path.userID": true
        },
        methodResponses: [
          {
            statusCode: "204",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "400",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "404",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          },
          {
            statusCode: "500",
            responseModels: {
              "application/json": apigateway.Model.EMPTY_MODEL
            }
          }
        ]
      })

    const usersTable = new Table(this, "users", {
      tableName: "users",
      partitionKey: {
        name: "id",
        type: AttributeType.STRING
      },
      removalPolicy: RemovalPolicy.DESTROY
    })
    usersTable.grantFullAccess(listUsers)
    usersTable.grantFullAccess(createUser)
    usersTable.grantFullAccess(getUser)
    usersTable.grantFullAccess(deleteUser)

    // The code that defines your stack goes here

    // example resource
    // const queue = new sqs.Queue(this, 'GoLambdaApiSampleQueue', {
    //   visibilityTimeout: cdk.Duration.seconds(300)
    // });
  }

  private get context(): StageContext {
    const env: string = this.node.tryGetContext("env");
    const context: StageContext = this.node.tryGetContext(env);
    context.ses = {
      fromEmail: process.env["SES_FROM_EMAIL"] || ""
    }

    return context
  }
}
