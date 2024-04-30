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
        data = response.json()
        # Render the template with JSON data
        return render_template('key_generation.html', data=data)
    else:
        return "Failed to fetch data from the Go endpoint", 500

if __name__ == '__main__':
    app.run(debug=True)
