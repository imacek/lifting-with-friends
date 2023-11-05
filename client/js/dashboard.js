/* globals Chart:false */

(() => {
  'use strict'

  fetch('/api/data', {
    method: 'GET',
    headers: {
      'Accept': 'application/json',
    },
  })
    .then(response => response.json())
    .then(response => {
      new Chart(document.getElementById('exerciseChart'), {
        type: 'line',
        data: {
          datasets: [
            { label: 'ivan', data: response['ivan']['Deadlift (Barbell)'] },
            { label: 'vinko',  data: response['vinko']['Deadlift (Barbell)'] },
            { label: 'linda', data: response['linda']['Deadlift (Barbell)'] },
            { label: 'yomach', data: response['yomach']['Deadlift (Barbell)'] },
          ]
        },
        options: {
          scales: {
            x: {
              type: 'time',
            }
          },
          parsing: {
            xAxisKey: 'timestamp',
            yAxisKey: 'maxOneRepMax'
          },
          plugins: {
            legend: {
              display: true
            },
            tooltip: {
              boxPadding: 3
            }
          }
        }
      })

    })

})()
