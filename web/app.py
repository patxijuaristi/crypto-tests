from flask import Flask, render_template
import requests
import database

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/real-time')
def real_time():
    data = {
        'page_title' : '📊 Real time testing data',
        'grafana_url' : 'http://localhost:3000/public-dashboards/d9399f5b4d704749a1eb0fa30d668fe4'
    }
    return render_template('grafana.html', data=data)

@app.route('/data-key-generation')
def data_key_generation():
    data = {
        'page_title' : '🔑 Key generation historical data',
        'grafana_url' : 'http://localhost:3000/public-dashboards/f39b8f60f2f64c1eaaa38985e343fda4'
    }
    return render_template('grafana.html', data=data)

@app.route('/test-key-generation')
def generate_key():
    # Make a GET request to the Go endpoint
    response = requests.get('http://localhost:8080/generateKey')
    
    # Check if the request was successful
    if response.status_code == 200:
        # Parse the JSON response
        test_result = response.json()
        test_result = set_bg_color(test_result)
        
        data = {
            'page_title' : '🔑 Key generation',
            'test_result' : test_result
        }
        # Render the template with JSON data
        return render_template('time_memory_results.html', data=data)
    else:
        return "Failed to fetch data from the Go endpoint", 500

@app.route('/test-signature-generation')
def generate_signature():
    # Make a GET request to the Go endpoint
    response = requests.get('http://localhost:8080/generateSignature')
    
    # Check if the request was successful
    if response.status_code == 200:
        # Parse the JSON response
        test_result = response.json()
        test_result = set_bg_color(test_result)
        
        data = {
            'page_title' : '✍ Signature generation',
            'test_result' : test_result
        }
        # Render the template with JSON data
        return render_template('time_memory_results.html', data=data)
    else:
        return "Failed to fetch data from the Go endpoint", 500

@app.route('/test-signature-verification')
def signature_verification():
    # Make a GET request to the Go endpoint
    response = requests.get('http://localhost:8080/verifySignature')
    
    # Check if the request was successful
    if response.status_code == 200:
        # Parse the JSON response
        test_result = response.json()
        test_result = set_bg_color(test_result)
        
        data = {
            'page_title' : '📝 Signature verification',
            'test_result' : test_result
        }
        # Render the template with JSON data
        return render_template('time_memory_results.html', data=data)
    else:
        return "Failed to fetch data from the Go endpoint", 500

@app.route('/signature-key-sizes')
def signature_key_sizes():
    data = {
        'page_title' : '📊 Key and Signature sizes',
        'grafana_url' : 'http://localhost:3000/public-dashboards/25123e0d5ad3464f88cd70ebf5bae07d'
    }
    return render_template('grafana.html', data=data)

def set_bg_color(data):
    for item in data:
        if 'Dilithium' in item['algorithm']:
            item['bg_color'] = '#fcffc3'
        elif 'SPHINCS+' in item['algorithm']:
            item['bg_color'] = '#aed4e9'
        elif 'Falcon' in item['algorithm']:
            item['bg_color'] = '#dddddd'
        else:
            item['bg_color'] = '#ffd3fc'
    return data

if __name__ == '__main__':
    app.run(debug=True)
