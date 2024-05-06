from flask import Flask, jsonify, render_template
import requests

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/key-generation')
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

@app.route('/signature-generation')
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

@app.route('/signature-verification')
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
    # Make a GET request to the Go endpoint
    response = requests.get('http://localhost:8080/keySignatureSizes')
    
    # Check if the request was successful
    if response.status_code == 200:
        # Parse the JSON response
        test_result = response.json()
        test_result = set_bg_color(test_result)
        
        data = {
            'page_title' : '📊 Key and Signature sizes',
            'test_result' : test_result
        }
        # Render the template with JSON data
        return render_template('sizes_results.html', data=data)
    else:
        return "Failed to fetch data from the Go endpoint", 500

def set_bg_color(data):
    for item in data:
        if 'Dilithium' in item['algorithm']:
            item['bg_color'] = '#a6d5a6'
        elif 'SPHINCS+' in item['algorithm']:
            item['bg_color'] = '#aed4e9'
        elif 'Falcon' in item['algorithm']:
            item['bg_color'] = '#fcffc3'
        else:
            item['bg_color'] = '#ffc3c3'
    return data

if __name__ == '__main__':
    app.run(debug=True)
