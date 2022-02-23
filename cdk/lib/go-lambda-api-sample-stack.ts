import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
// import * as sqs from 'aws-cdk-lib/aws-sqs';
import { aws_apigateway as apigateway } from "aws-cdk-lib";
import { Code, Function, Runtime } from 'aws-cdk-lib/aws-lambda';
import { StageContext } from './context';
export class GoLambdaApiSampleStack extends Stack {
  constructor(scope: Construct, id: string, props?: StackProps) {
    super(scope, id, props);

    const env: string = this.node.tryGetContext("env");
    const context: StageContext = this.node.tryGetContext(env);

    /** functions **/
    const listUsers = new Function(this, "listUsers", {
      description: "list users",
      code: Code.fromAsset('../bin/zip/list-users.zip'),
      handler: 'list-users',
      runtime: Runtime.GO_1_X,
    })
    const getUser = new Function(this, "getUser", {
      description: "get user by id",
      code: Code.fromAsset('../bin/zip/get-user.zip'),
      handler: 'get-user',
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
            statusCode: "404",
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
    const user = users.addResource("{user_id}")
    user.addMethod("GET",
      new apigateway.LambdaIntegration(getUser, {
        proxy: true,
        passthroughBehavior: apigateway.PassthroughBehavior.WHEN_NO_TEMPLATES,
        requestTemplates: {
          "application/json": "{ \"userID\": \"$input.params('user_id')\" }"
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
    // The code that defines your stack goes here

    // example resource
    // const queue = new sqs.Queue(this, 'GoLambdaApiSampleQueue', {
    //   visibilityTimeout: cdk.Duration.seconds(300)
    // });
  }
}
