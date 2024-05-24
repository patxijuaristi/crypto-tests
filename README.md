Post Quantum Cryptography algorithms 💻
================================================================

This repository is composed by a main program that performs tests to compare the performance of post quantum cryptography methods on a blockchain network, in this case in Geth (go-ethereum) blockchain. The [modified Go-Ethereum](https://github.com/patxijuaristi/private-ethereum-pq) calls the API that performs the tests by comparing various cryptographic algorithms and stores the results in a Sqlite database. This code is developed in GO, and used algorithms are: 

- ECDSA (Default algorihtm used in Geth)
- SPHINCS+
- Crystals Dilithium
- Falcon

These results are displayed on a web page made with Flask (found in the */web* folder) using graphs made in Grafana (found exported in the */grafana* folder).

## Installation

The instructions for installing the test environment on a new system are detailed below.

First it is necessary to install and run the private Blockchain network from my other repository:

https://github.com/patxijuaristi/private-ethereum-pq

Once we have the network ready, we can start with the installation of the required systems.

### GO program

In order to execute the GO code, which is the main program, it will be necessary to install GO itself and the necessary libraries listed in the `go.mod` file. Once we have the environment ready, we can execute the code in two ways: in simulation mode, and API endpoint mode. In the simulation mode, there is a test done every minute with a randomly generated hash. In API mode, the code runs the tests with the hash values received from the private Ethereum blockchain, which calls the API enpodints including that value. The tests are executed asynchronously using a queuing system, because the API receives more than one call in the time it needs to perform a single test.

To execute the code in API mode:

```
go run . 1
```

To execute the code in simulation mode:

```
go run . 2
```

### Python

The website is developed in Python, using Flask, which is a framework for web development in that language. Therefore, the first step is to install Python.

Once Python is installed, it will be necessary to install `virtualenv` in order to create the virtual environment that we will create to have all the necessary dependencies for the execution of the web.

```
pip install virtualenv
```

With `virtualenv`, we are going to create the virtual environment, which will be called `env`.

```
virtualenv env
```

With the environment created, we activate it and enter it.

```
.\env\Scripts\activate
```

And we install all the necessary libraries with the file `requirements.txt`.

```
pip install -r .\requirements.txt
```

### Grafana

Next, we need to install Grafana on our system, and set the data source from which it will read the information to display the results. This will be the path to the `crypto_tests.db` file which is automatically created on the first run of the main program.

![Grafana Data Sources](https://github.com/patxijuaristi/crypto-tests/blob/master/images/grafana-data-source.png)

Then we need to import all the dashboards that are exported in JSON format in the */grafana* folder. The result should be the next one:

![Grafana Dashboards](https://github.com/patxijuaristi/crypto-tests/blob/master/images/grafana-dashboards.png)

Then to visualize them in the local web pages, we need to make them public and change the URL-s. We have to go one by one in the dashboards, and get the public address, copy it and set it in the `app.py` file, in the corresponding view.

![Grafana Share](https://github.com/patxijuaristi/crypto-tests/blob/master/images/grafana-share.png)

```
'grafana_url' : 'http://localhost:3000/public-dashboards/xxxxxxxxxxxxxxxxxxxxxxxxx'
```