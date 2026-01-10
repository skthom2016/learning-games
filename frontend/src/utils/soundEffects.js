// Sound effects utility using Web Audio API
// This avoids the need for external audio files

const audioContext = new (window.AudioContext || window.webkitAudioContext)();

// Play a cheerful sound for correct answers
export const playCorrectSound = () => {
  try {
    const oscillator = audioContext.createOscillator();
    const gainNode = audioContext.createGain();

    oscillator.connect(gainNode);
    gainNode.connect(audioContext.destination);

    // Create a cheerful ascending arpeggio
    const notes = [523.25, 659.25, 783.99, 1046.50]; // C5, E5, G5, C6 (C major chord)
    let noteIndex = 0;

    const playNote = () => {
      if (noteIndex >= notes.length) {
        oscillator.disconnect();
        return;
      }

      const now = audioContext.currentTime;
      oscillator.frequency.setValueAtTime(notes[noteIndex], now);
      oscillator.type = 'sine';

      gainNode.gain.setValueAtTime(0.3, now);
      gainNode.gain.exponentialRampToValueAtTime(0.01, now + 0.2);

      noteIndex++;
      setTimeout(playNote, 100);
    };

    oscillator.start();
    playNote();

    setTimeout(() => {
      oscillator.stop();
      oscillator.disconnect();
    }, 500);
  } catch (error) {
    console.error('Error playing correct sound:', error);
  }
};

// Play a gentle "oh-oh" sound for incorrect answers
export const playIncorrectSound = () => {
  try {
    const oscillator = audioContext.createOscillator();
    const gainNode = audioContext.createGain();

    oscillator.connect(gainNode);
    gainNode.connect(audioContext.destination);

    const now = audioContext.currentTime;

    // Create a descending tone (gentle "oh-oh" effect)
    oscillator.frequency.setValueAtTime(392.00, now); // G4
    oscillator.frequency.exponentialRampToValueAtTime(329.63, now + 0.15); // Drop to E4
    oscillator.type = 'sine';

    gainNode.gain.setValueAtTime(0.2, now);
    gainNode.gain.exponentialRampToValueAtTime(0.01, now + 0.3);

    oscillator.start(now);
    oscillator.stop(now + 0.3);

    setTimeout(() => {
      oscillator.disconnect();
      gainNode.disconnect();
    }, 400);
  } catch (error) {
    console.error('Error playing incorrect sound:', error);
  }
};

// Play a celebration fanfare (multiple notes)
export const playCelebrationSound = () => {
  try {
    const playFanfare = () => {
      const notes = [523.25, 659.25, 783.99, 1046.50, 1318.51]; // C5, E5, G5, C6, E6
      notes.forEach((freq, index) => {
        setTimeout(() => {
          const osc = audioContext.createOscillator();
          const gain = audioContext.createGain();

          osc.connect(gain);
          gain.connect(audioContext.destination);

          const now = audioContext.currentTime;
          osc.frequency.setValueAtTime(freq, now);
          osc.type = 'triangle';

          gain.gain.setValueAtTime(0.2, now);
          gain.gain.exponentialRampToValueAtTime(0.01, now + 0.3);

          osc.start(now);
          osc.stop(now + 0.3);

          setTimeout(() => {
            osc.disconnect();
            gain.disconnect();
          }, 400);
        }, index * 80);
      });
    };

    playFanfare();
  } catch (error) {
    console.error('Error playing celebration sound:', error);
  }
};
