# split CPU cores, based on share of each service, create multiple containers. at most 1 CPU core per container

# https://github.cpom/ppgaluzio/MOBOpt
# @article{GALUZIO2020100520,
# title = "MOBOpt — multi-objective Bayesian optimization",
# journal = "SoftwareX",
# volume = "12",
# pages = "100520",
# year = "2020",
# issn = "2352-7110",
# doi = "https://doi.org/10.1016/j.softx.2020.100520",
# url = "http://www.sciencedirect.com/science/article/pii/S2352711020300911",
# author = "Paulo Paneque Galuzio and Emerson Hochsteiner [de Vasconcelos Segundo] and Leandro dos Santos Coelho and Viviana Cocco Mariani"
# }
import json
import mobopt as mo
import numpy as np
import sys
core_count = 20000 # make sure this is the same as what you have in config file of acfg as ConfigurationValidation.TotalCPU
cache = {}

import time
def objective(x):

    s = sum(x)
    g = core_count * (x[0] / s)
    a = core_count * (x[1] / s)
    b = core_count * (x[2] / s)
    # this configuration would use a + b + g cores. x[3] is the amount which is not being used

    config = {
        'gateway': {
            "cpu_count": np.round(g, 0),
        },
        'auth': {
            "cpu_count": np.round(a, 0),
        },
        'books': {
            "cpu_count": np.round(b, 0),
        },
    }

    key = json.dumps(config, sort_keys=True)
    if key in cache:
        return cache[key]

    print(json.dumps(config), flush=True)
    line = "default"
    for line in sys.stdin:
        data = json.loads(line.strip())
        break
    with open("/home/vahid/Desktop/log.python.mobo", "w+") as f:
        f.write(str(x))
    f.close()

    SLA_target = 250
    respones_times = [0] * 3
    for i in range(3):
        if data['feedbacks'][i] > SLA_target:
            respones_times[i] = data['feedbacks'][i] - SLA_target
    
    res = [respones_times[0],respones_times[1],respones_times[2],a+b+g]
    cache[key] = np.array(res)

    return np.array(np.array(res))

PB = np.asarray([
    [0.06, 0.94],
    [0.06, 0.94],
    [0.06, 0.94],
    [0.06, 0.94]
])
NParam = PB.shape[0]

Optimizer = mo.MOBayesianOpt(target=objective,
                             NObj=4,
                             pbounds=PB,
                             verbose=False,
                             max_or_min='min', # whether the optimization problem is a maximization problem ('max'), or a minimization one ('min')
                             RandomSeed=10)
Optimizer.initialize(init_points=5) 
# there is no minimize function. maximize() starts optimization. Performs minimizing or maximizing based on max_or_min
front, pop = Optimizer.maximize(n_iter=20,
                                prob=0.1,
                                ReduceProb=False,
                                q=0.5)
print('done')