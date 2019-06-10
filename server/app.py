from pymongo import MongoClient
from concurrent import futures
import time
import logging

import grpc

import crud_pb2
import crud_pb2_grpc
import pprint

_ONE_DAY_IN_SECONDS = 60 * 60 * 24

class CRUD(crud_pb2_grpc.CRUDServicer):

    def Insert(self, request, context):
        client = MongoClient()
        db = client['crud-grpc']
        employee = db.employees

        data_insert= {
            'name' : request.name,
            'city' : request.city,
        }

        employee.insert_one(data_insert)

        return crud_pb2.StatusResponse(message='Success Insert Data With Name: ' + request.name)


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    crud_pb2_grpc.add_CRUDServicer_to_server(CRUD(), server)
    server.add_insecure_port('[::]:55551')
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
