<template>
  <div class="management-buttons">
    <button @click="startRecording" class="record-button">Start Recording</button>
    <button @click="stopRecording" class="stop-button">Stop Recording</button>

    <p v-if="statusMessage">{{ statusMessage }}</p>

  </div>
  <div>
    <h1>Video Files</h1>
    <table>
      <thead>
      <tr>
        <th>#</th>
        <th>File Name</th>
        <th>Actions</th>
      </tr>
      </thead>
      <tbody>
      <tr v-for="(video, index) in videos" :key="index">
        <td>{{ index + 1 }}</td>
        <td>{{ video.name }}</td>
        <td>
          <button>
            <img
                v-if="video.thumbnail" :src="`data:image/png;base64,${video.thumbnail}`"
                alt="Thumbnail" style="width: 150px; height: auto;"
                @click="playVideo(video.name)">
          </button>
        </td>
      </tr>
      </tbody>
    </table>
    <video v-if="videoBlobUrl" autoplay controls style="width: 100%; margin-top: 20px;">
      <source :src="videoBlobUrl" type="video/mp4">
      Your browser does not support the video tag.
    </video>

  </div>
</template>

<script>
export default {
  data() {
    return {
      videos: [],
      videoBlobUrl: null,
      statusMessage: '',
    };
  },
  mounted() {
    this.fetchVideos();
  },
  methods: {
    fetchVideos() {
      fetch('http://localhost:8080/api/list_videos') // URL API для получения списка видео
          .then(response => response.json())
          .then(data => {
            console.log("data: " + data.videos)
            this.videos = data.videos;
          })
          .catch(error => console.error('Error fetching videos list:', error));
    },
    playVideo(video) {
      console.log("playing video: " + JSON.stringify(video))
      this.videoBlobUrl = ""
      fetch(`http://localhost:8080/api/${video}`, {
        method: "POST",
      })
          .then(response => {
            if (!response.ok) {
              throw new Error("Failed to fetch videos")
            }
            return response.blob()
          })
          .then(blob => {
            this.videoBlobUrl = URL.createObjectURL(blob)
          })
          .catch(error => console.error('Error playing videos: ', error))
    },
    startRecording() {
      fetch('http://localhost:8080/api/start_multiply_recording', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
      })
          .then(response => response.json())
          .then(data => {
            if (data.error) {
              this.statusMessage = `Error: ${data.error}`;
            } else {
              this.statusMessage = data.message;
            }
          })
          .catch(error => {
            console.error('Error starting recording:', error);
            this.statusMessage = 'Error starting recording';
          });
    },
    // Остановить запись
    stopRecording() {
      fetch('http://localhost:8080/api/stop_multiply_recording', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
      })
          .then(response => response.json())
          .then(data => {
            if (data.error) {
              this.statusMessage = `Error: ${data.error}`;
            } else {
              this.statusMessage = data.message;
            }
          })
          .catch(error => {
            console.error('Error stopping recording:', error);
            this.statusMessage = 'Error stopping recording';
          });
    }

  }
};
</script>

<style scoped>
table {
  width: 100%;
  border-collapse: collapse;
}

th, td {
  padding: 10px;
  text-align: left;
  border: 1px solid #ddd;
}

th {
  background-color: #f4f4f4;
}

.management-buttons {
  display: flex;
}

.record-button {
  background-color: green;
  margin: 10px;
}

.stop-button {
  background-color: red;
  margin: 10px;
}

button {
  padding: 5px 10px;
  color: #fff;
  background-color: #007bff;
  border: none;
  border-radius: 3px;
  cursor: pointer;
}

button:hover {
  background-color: #0056b3;
}


</style>