from pymongo import MongoClient
from concurrent import futures
import time
import logging
import json
from bson.objectid import ObjectId

import grpc

import crud_pb2
import crud_pb2_grpc
import pprint

_ONE_DAY_IN_SECONDS = 60 * 60 * 24

class CRUD(crud_pb2_grpc.CRUDServicer):

    def Connection(self):
        client = MongoClient('mongo', 27017)
        db = client['crud-grpc']
        employee = db.employees

        return employee

    def Insert(self, request, context):
        employee = self.Connection()

        data_insert= {
            'name' : request.name,
            'city' : request.city,
        }

        employee.insert_one(data_insert)
        return crud_pb2.StatusResponse(message='Success Insert Data With Name: ' + request.name)

    def List(self, request, context):
        employee = self.Connection()

        employees = employee.find({})
        
        for data in employees:
            if data is not None:
                listdata = crud_pb2.DataResponse(
                    id=str(data["_id"]),
                    name=data["name"],
                    city=data["city"]
                )
                yield listdata

    def Show(self, request, context):
        employee = self.Connection()

        employees = employee.find_one({"_id": ObjectId(request.id)})

        return crud_pb2.DataResponse(
            id=str(employees["_id"]),
            name=employees["name"],
            city=employees["city"]
        )

    def Update(self, request, context):
        employee = self.Connection()

        data_update = {
            'name': request.name,
            'city': request.city,
        }
        employee.replace_one({"_id": ObjectId(request.id)}, data_update)

        return crud_pb2.StatusResponse(message='Success Update Data With Name: ' + request.name)

    def Delete(self, request, context):
        employee = self.Connection()
        
        employee.delete_one({'_id': ObjectId(request.id)})

        return crud_pb2.StatusResponse(message="Success Delete Data With Id: " + request.id)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    crud_pb2_grpc.add_CRUDServicer_to_server(CRUD(), server)
    server.add_insecure_port('0.0.0.0:44441')
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    print('Starting Server...')
    logging.basicConfig()
    serve()

# results = courses.insert_many(arr_course)

# for object_id in results.inserted_ids:
#     print('Course Added. The course id is ' + str(object_id))

# courses.update({
#     'course' : 'MongoDB Tutorial'
# },
# {
#     '$set' : {
#         'course' : 'MongoDB Complete Guide'
#     }
# }, multi=True)

# print(courses.find({'author' : 'Ridwan'}).count())

# courses.delete_many({
#     'author' : 'Ridwan'
# })

# print(courses.find({'author' : 'Ridwan'}).count())

# courses = courses.find({'course' : 'MongoDB Complete Guide'})

# for course in courses:
#     pprint.pprint(course)

# result = courses.insert_one(course)

# if result.acknowledged:
#     print('Course Added. The course id is ' + str(result.inserted_id))

# print(list(courses.aggregate([
#     {
#         '$group': {
#             '_id': '$author',
#             'rangking' : {
#                 '$avg' : '$rating'
#             }
#         }
#     }
# ])))
