from itertools import combinations
# import networkx as nx
# import matplotlib.pyplot as plt
# import numpy as np
import random
import math


class OptimizationProblem:
    def __init__(self, vehicles, nodes, distance_matrix):
        self.vehicles = vehicles
        self.nodes = nodes
        self.distance_matrix = distance_matrix

        # Decision variables
        self.x_ijk = {
            (i, j, k): 0 for i in self.nodes for j in self.nodes for k in self.vehicles}
        self.y_ik = {(i, k): 0 for i in self.nodes for k in self.vehicles}


class Vehicle:
    def __init__(self, vehicle_id, fixed_cost, capacity, velocity, colony_id, variable_cost, gas_emissions):
        self.vehicle_id = vehicle_id
        self.fixed_cost = fixed_cost
        self.capacity = capacity
        self.velocity = velocity
        self.colony_id = colony_id
        self.variable_cost = variable_cost
        self.gas_emissions = gas_emissions


class DistanceMatrix:
    def __init__(self, num_nodes):
        self.num_nodes = num_nodes
        self.distances = [[0] * num_nodes for _ in range(num_nodes)]

    def set_distance(self, node_i, node_j, distance):
        self.distances[node_i][node_j] = distance
        self.distances[node_j][node_i] = distance

    def get_distance(self, node_i, node_j):
        return self.distances[node_i][node_j]


class Node:
    def __init__(self, node_id, location, start_time, end_time, service_time, demands):
        self.node_id = node_id
        self.location = location
        self.start_time = start_time  # Time window's lower bound
        self.end_time = end_time
        self.service_time = service_time
        self.waiting_time = 0
        self.demands = demands  # Dictionary of demand for each product type
        self.travel_times = {}  # Use Node objects as keys for travel_times


# ATTRIBUTES VALUES
# Create Nodes(node_id, location, start_time, end_time, service_time, demands)
nodes = [
    Node(0, (50, 50), 0, 24, 0.5, {}),
    Node(1, (70, 70), 8, 17, 0.5, 100),
    Node(2, (30, 40), 8, 17, 0.8, 100),
    Node(3, (20, 80), 8, 17, 1.2, 100),
    Node(4, (90, 20), 8, 17, 0.7, 100),
    Node(5, (40, 20), 8, 17, 1.0, 100),
    Node(6, (80, 30), 8, 17, 0.6, 100),
    Node(7, (60, 70), 8, 17, 0.9, 100),
    Node(8, (40, 60), 8, 17, 1.1, 100),
    Node(9, (65, 55), 8, 17, 0.5, 100)
]


# Create distance matrix
num_nodes = len(nodes)
distance_matrix = DistanceMatrix(num_nodes)

# Calculate distances between nodes and set them in the distance matrix
for i in range(num_nodes):
    for j in range(i + 1, num_nodes):
        distance = math.sqrt((nodes[i].location[0] - nodes[j].location[0]) ** 2 +
                             (nodes[i].location[1] - nodes[j].location[1]) ** 2)
        distance_matrix.set_distance(i, j, distance)


# Create Vehicle(self, vehicle_id, fixed_cost, capacity, velocity, colony_id, variable_cost, gas_emissions)
vehicles = [
    Vehicle(1, 100, 400, 50, 0, 5.5, 0.5),
    Vehicle(2, 150, 500, 60, 1, 7.5, 0.6),
    Vehicle(3, 250, 700, 70, 2, 8.7, 0.8)
]
# MATHEMATICAL MODEL (OBJECTIVE FUNCTIONS)

# Create an instance of OptimizationProblem
opt_problem = OptimizationProblem(vehicles, nodes, distance_matrix)

# Create a solution object based on the decision variables
solution = {
    'x_ijk': opt_problem.x_ijk,
    'y_ik': opt_problem.y_ik,
    'opt_problem': opt_problem
}


# Objective function 1: Calculate the total cost
def objective_function_1(x_ijk, opt_problem):
    distance_matrix = opt_problem.distance_matrix

    total_fixed_cost = sum(
        k.fixed_cost * x_ijk[0, j.node_id, k.vehicle_id]
        for j in opt_problem.nodes
        for k in opt_problem.vehicles
        if (0, j.node_id, k.vehicle_id) in x_ijk
    )

    total_variable_cost = sum(
        distance_matrix.get_distance(
            i.node_id, j.node_id) * k.variable_cost * x_ijk[i.node_id, j.node_id, k.vehicle_id]
        for k in opt_problem.vehicles
        for i in opt_problem.nodes
        for j in opt_problem.nodes
        if (i.node_id, j.node_id, k.vehicle_id) in x_ijk
    )

    # Calculate the total cost
    total_cost = total_fixed_cost + total_variable_cost
    return total_cost


# Objective function 3: Calculate the total time-related costs
def objective_function_3(x_ijk, y_ik, opt_problem):
    distance_matrix = opt_problem.distance_matrix

    total_travel_time = sum(
        distance_matrix.get_distance(
            i.node_id, j.node_id) / k.velocity * x_ijk[i.node_id, j.node_id, k.vehicle_id]
        for k in opt_problem.vehicles
        for i in opt_problem.nodes
        for j in opt_problem.nodes
        if (i.node_id, j.node_id, k.vehicle_id) in x_ijk
    )

    total_service_time = sum(
        i.service_time * y_ik[i.node_id, k.vehicle_id]
        for i in opt_problem.nodes
        for k in opt_problem.vehicles
        if (i.node_id, k.vehicle_id) in y_ik
    )

    total_waiting_time = sum(
        i.waiting_time * y_ik[i.node_id, k.vehicle_id]
        for i in opt_problem.nodes
        for k in opt_problem.vehicles
        if (i.node_id, k.vehicle_id) in y_ik
    )

    total_time = total_travel_time + total_service_time + total_waiting_time

    return total_time


# Objective function 4: Calculate the total gas emissions
def objective_function_4(x_ijk, opt_problem):
    distance_matrix = opt_problem.distance_matrix

    total_gas_emissions = sum(
        distance_matrix.get_distance(
            i.node_id, j.node_id) * k.gas_emissions * x_ijk[i.node_id, j.node_id, k.vehicle_id]
        for k in opt_problem.vehicles
        for i in opt_problem.nodes
        for j in opt_problem.nodes
        if (i.node_id, j.node_id, k.vehicle_id) in x_ijk
    )

    return total_gas_emissions
# MATHEMATICAL MODEL (CONSTRAINTS)


# Constraint 1: Each vehicle can leave at most once for each trip
def single_departure_constraint(x_ijk, opt_problem):
    sum_leave_once = sum(x_ijk[0, j.node_id, k.vehicle_id]
                         for j in opt_problem.nodes
                         for k in opt_problem.vehicles
                         if (0, j.node_id, k.vehicle_id) in x_ijk
                         )
    if sum_leave_once > 1:
        return False  # Constraint violation detected

    return True  # All constraints satisfied


# Constraint 2: Preservation flow constraint
def preservation_flow_constraint(x_ijk, y_ik, opt_problem):
    sum_successors = sum(x_ijk[i.node_id, j.node_id, k.vehicle_id]
                         for i in opt_problem.nodes
                         for j in opt_problem.nodes if j.node_id > i.node_id
                         for k in opt_problem.vehicles
                         if (i.node_id, j.node_id, k.vehicle_id) in x_ijk
                         )

    sum_predecessors = sum(x_ijk[j.node_id, i.node_id, k.vehicle_id]
                           for i in opt_problem.nodes
                           for j in opt_problem.nodes if j.node_id < i.node_id
                           for k in opt_problem.vehicles
                           if (i.node_id, j.node_id, k.vehicle_id) in x_ijk
                           )
    lhs = sum_successors + sum_predecessors
    rhs = 2 * sum(y_ik[i.node_id, k.vehicle_id]
                  for i in opt_problem.nodes
                  for k in opt_problem.vehicles
                  if (i.node_id, k.vehicle_id) in y_ik
                  )
    if lhs != rhs:
        return False  # Constraint violation detected

    return True  # All constraints satisfied


# Constraint 3: Subtour elimination constraint
def subtour_elimination_constraint(x_ijk, opt_problem):
    for subset_size in range(2, len(opt_problem.nodes) + 1):
        for subset in combinations(opt_problem.nodes, subset_size):
            subset_sum = sum(x_ijk[i.node_id, j.node_id, k.vehicle_id]
                             for i in subset
                             for j in subset if i != j
                             for k in opt_problem.vehicles
                             if (i.node_id, j.node_id, k.vehicle_id) in x_ijk
                             )
            if subset_sum > subset_size - 1:
                return False  # Constraint violation, solution is infeasible

    return True  # All constraints are satisfied, solution is feasible


# Constraint 5: Max delivery quantity constraint
def max_delivery_quantity_constraint(y_ik, opt_problem):
    lhs = sum(i.demands * y_ik[i.node_id, k.vehicle_id]
              for k in opt_problem.vehicles
              for i in opt_problem.nodes
              if (i.node_id, k.vehicle_id) in y_ik
              )

    rhs = sum(k.capacity
              for k in opt_problem.vehicles
              for i in opt_problem.nodes
              if (i.node_id, k.vehicle_id) in y_ik
              )

    if lhs > rhs:
        return False  # Constraint violation

    return True  # All constraints satisfied


# Constraint 12: Binary variable constraints
def binary_variable_constraints(x_ijk, y_ik, opt_problem):

    # Check x_ijk binary constraints
    for i in opt_problem.nodes:
        for j in opt_problem.nodes:
            for k in opt_problem.vehicles:
                if (i.node_id, j.node_id, k.vehicle_id) in x_ijk:
                    if x_ijk[i.node_id, j.node_id, k.vehicle_id] != 0 and x_ijk[i.node_id, j.node_id, k.vehicle_id] != 1:
                        return False  # Constraint violation

    # Check y_ik binary constraints
    for i in opt_problem.nodes:
        for k in opt_problem.vehicles:
            if (i.node_id, k.vehicle_id) in y_ik:
                if y_ik[i.node_id, k.vehicle_id] != 0 and y_ik[i.node_id, k.vehicle_id] != 1:
                    return False  # Constraint violation

    return True  # Constraints satisfied
# ANT COLONY OPTIMIZATION


# Define the evaluate_solution() function

def evaluate_solution(x_ijk, y_ik, opt_problem):
    # Calculate the fitness value or objective score based on the evaluation criteria
    obj1 = objective_function_1(x_ijk, opt_problem)
    obj3 = objective_function_3(x_ijk, y_ik, opt_problem)
    obj4 = objective_function_4(x_ijk, opt_problem)

    # Evaluate the problem constraints
    constraints_satisfied = all([
        single_departure_constraint(x_ijk, opt_problem),
        preservation_flow_constraint(x_ijk, y_ik, opt_problem),
        subtour_elimination_constraint(x_ijk, opt_problem),
        max_delivery_quantity_constraint(y_ik, opt_problem),
        binary_variable_constraints(x_ijk, y_ik, opt_problem)
    ])

    # Assign a large penalty if any of the constraints are violated
    if not constraints_satisfied:
        penalty = 1000000  # Example: Use a large penalty value
        obj1 += penalty
        obj3 += penalty
        obj4 += penalty

    # Combine the objective values into a single fitness value or objective score
    # Example: Use a tuple to represent the multi-objective fitness
    fitness = (obj1, obj3, obj4)

    return fitness


# Evaluate the solution using the objective functions
fitness = evaluate_solution(
    solution['x_ijk'],
    solution['y_ik'],
    solution['opt_problem']
)

# Define ACO parameters
num_colonies = 3
total_iterations = 1
pheromone_evaporation_rate = 1
alpha = 0.1
beta = 2.0

# Create multiple colonies for each type of vehicle in the depot
colonies = []
for colony_id in range(num_colonies):
    colony = {
        'colony_id': colony_id,
        'ants': [],
        'taboo_list': set(),
        # Add other colony-specific parameters
    }
    colonies.append(colony)

# Initialize pheromone matrix
num_nodes = len(opt_problem.nodes)
pheromone_matrix = [[1.0] * num_nodes for _ in range(num_nodes)]

# Main ACO algorithm loop
best_solution = None
best_fitness = None
all_solution = []
for iteration in range(total_iterations):
    print("ITERACION-------------------------------------------------------------------------------------------------")
    # Construct ant solutions for each colony
    for colony in colonies:
        colony_vehicles = [
            vehicle for vehicle in opt_problem.vehicles if vehicle.colony_id == colony['colony_id']]
        colony_best_solution = None
        colony_best_fitness = None
        taboo_list = colony['taboo_list'].copy()  # Define taboo_list here

        for ant_id in range(len(colony_vehicles)):
            print("VEHICULO------------------------------------------------")
            ant = {
                'position': 0,  # Start at the depot
                # List to store visited nodes (initialize with depot)
                'solution': [0],
                # 'capacity': [0] * len(colony_vehicles),  # Initialize capacity for each vehicle in the colony
                'capacity': 200,
                # Add other ant-related parameters
            }

            # Visit customers until capacity is exhausted
            while ant['capacity'] > 0:
                # Compute probabilities for selecting the next node
                probabilities = []
                for node_id in range(1, len(opt_problem.nodes)):
                    if node_id not in taboo_list:
                        pheromone = pheromone_matrix[ant['position']][node_id]
                        distance = opt_problem.distance_matrix.get_distance(
                            ant['position'], node_id)
                        attractiveness = 1 / distance  # Heuristic information
                        probability = pheromone ** alpha * attractiveness ** beta
                        probabilities.append((node_id, probability))
                # print(probabilities)
                # print("--random_probabilities----")
                random.shuffle(probabilities)
                print(probabilities)
               # print("-random_value----")
                # Select the next node based on the probabilities
                total_probability = sum(prob for _, prob in probabilities)

                random_value = random.uniform(0, total_probability)
                # print("---total_probability---")
                # print(total_probability)
                # print("--random_value----")
                # print(random_value)
                # print("------")
                cumulative_probability = 0
                next_node = None
                # print(probabilities)
                for node_id, probability in probabilities:
                    cumulative_probability += probability
                    # print("---cumulative_probability---")
                   # print(cumulative_probability)
                    # print("---for---")
                    # print(probability)
                   # print("---nodo---")
                   # print(node_id)
                    if cumulative_probability >= random_value:
                        next_node = node_id
                        break

                if next_node != None:

                 # Update the taboo list and subtract the demand from the vehicle's capacity
                    taboo_list.add(next_node)
                    # print(next_node)
                    node = opt_problem.nodes[next_node]
                    # print("---demands---")
                    ant['capacity'] -= node.demands
                    ant['solution'].append(next_node)

                print(ant['capacity'])
                # Update the ant's position and solution

                if ant['capacity'] <= 0:
                    ant['solution'].append(0)
                    ant['capacity'] = 0

            # Calculate fitness for the ant's solution
            print(ant['solution'])
            ant_fitness = evaluate_solution(
                solution['x_ijk'],
                solution['y_ik'],
                solution['opt_problem']
            )

            # Update the best solution found so far in this colony
            if colony_best_fitness is None or ant_fitness < colony_best_fitness:
                colony_best_fitness = ant_fitness
                colony_best_solution = ant['solution']

            # Add the ant's solution to the colony
            colony['ants'].append(ant)

    # Update pheromone trails based on the best ant's solution in this colony
    if colony_best_solution is not None:
        for i in range(len(colony_best_solution) - 1):
            current_node = colony_best_solution[i]
            next_node = colony_best_solution[i + 1]
            # Update pheromone value
            pheromone_matrix[current_node][next_node] += 1

    # Update the colony's taboo list
    colony['taboo_list'].update(taboo_list)

    # Select the best solution among all colonies and update the global best
    for colony in colonies:
        for ant in colony['ants']:
            all_solution.append(ant['solution'])
            ant_fitness = evaluate_solution(
                solution['x_ijk'],
                solution['y_ik'],
                solution['opt_problem']
            )
            if best_fitness is None or ant_fitness < best_fitness:
                best_fitness = ant_fitness
                best_solution = ant['solution']

    # Clear ants in each colony for the next iteration
    for colony in colonies:
        colony['ants'] = []

# Output the best solution
print("All Solution:", all_solution)
print("Best Solution:", best_solution)
print("Solution Fitness:", best_fitness)
