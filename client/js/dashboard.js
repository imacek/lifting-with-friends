/* globals Chart:false */

(() => {
  'use strict'

  let data = null;
  let selectedMetric = null;

  const metricSelectBtnElems = document.getElementsByClassName('metric-select-btn');

  function handleMetricSelectButtonClicked(clickedBtnElem) {
    selectedMetric = clickedBtnElem.getAttribute('metric')

    Array.from(metricSelectBtnElems).forEach(element => {
      if (element.getAttribute('metric') === selectedMetric) {
        element.classList.add('btn-primary')
        element.classList.remove('btn-outline-secondary')
      } else {
        element.classList.add('btn-outline-secondary')
        element.classList.remove('btn-primary')
      }
    });

    reloadCharts()
  }

  handleMetricSelectButtonClicked(Array.from(metricSelectBtnElems)[0]);

  Array.from(metricSelectBtnElems).forEach(element => {
    element.onclick = event => handleMetricSelectButtonClicked(event.target)
  })

  function deleteCharts() {
    Array.from(document.getElementsByClassName('exercise-chart')).forEach(element => {
      element.remove();
    });
  }

  function showChart(exercise, metric, data) {

    const canvasElem = document.createElement("canvas");
    canvasElem.className = 'exercise-chart my-4 w-100'
    canvasElem.width = 900
    canvasElem.height = 380

    const mainElem = document.getElementsByTagName('main')[0];
    mainElem.insertAdjacentElement('beforeend', canvasElem)

    const datasets = Object.entries(data).map(([user, userData]) => {
      return { label: user, data: userData[exercise] }
    })

    new Chart(canvasElem, {
      type: 'line',
      data: {
        datasets: datasets
      },
      options: {
        scales: {
          x: {
            type: 'time',
          }
        },
        parsing: {
          xAxisKey: 'timestamp',
          yAxisKey: metric
        },
        plugins: {
          title: {
            display: true,
            text: exercise,
          },
          legend: {
            display: true
          }
        }
      }
    })
  }

  function reloadCharts() {
    if (!data) {
      return
    }

    deleteCharts()
    showChart('Squat (Barbell)', selectedMetric, data)
    showChart('Deadlift (Barbell)', selectedMetric, data)
    showChart('Bench Press (Barbell)', selectedMetric, data)
    showChart('Overhead Press (Barbell)', selectedMetric, data)
    showChart('Bent Over Row (Barbell)', selectedMetric, data)
  }

  async function uploadExerciseCsv()
  {
    let userElem = document.getElementById("user-input")
    let fileElem = document.getElementById("file-input")

    let formData = new FormData();
    formData.append("user", userElem.value);
    formData.append("file", fileElem.files[0]);

    const abortCtrl = new AbortController()
    setTimeout(() => abortCtrl.abort(), 5000);

    try {
      let r = await fetch('/api/upload', { method: "POST", body: formData, signal: abortCtrl.signal });
      console.log('HTTP response code:',r.status);
    } catch(e) {
      console.log('Huston we have problem...:', e);
    }
  }

  document.getElementById("upload-button").onclick = uploadExerciseCsv;


  fetch('/api/data', {
    method: 'GET',
    headers: {
      'Accept': 'application/json',
    },
  })
    .then(response => response.json())
    .then(response => {
      data = response
      reloadCharts()
    })

})()
