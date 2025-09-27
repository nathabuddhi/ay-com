!pip install transformers torchaudio
!pip install fsspec==2023.9.2
!pip install -U datasets

from datasets import load_dataset, Audio
import os

dataset = load_dataset("PolyAI/minds14", "en-US", split="train", trust_remote_code=True)
splits = dataset.train_test_split(test_size=0.2)
train_dataset = splits["train"]
test_dataset = splits["test"]

import torch
import torchaudio
import numpy as np

def resample_audio(sample):
    original_sr = sample["audio"]["sampling_rate"]
    target_sr = 16000
    if original_sr != target_sr:
        waveform = torch.tensor(sample["audio"]["array"], dtype=torch.float32).unsqueeze(0)
        resampler = torchaudio.transforms.Resample(original_sr, target_sr)
        waveform = resampler(waveform)
        sample["audio"]["array"] = waveform.squeeze(0).numpy()
        sample["audio"]["sampling_rate"] = target_sr

    if np.isnan(sample["audio"]["array"]).any() or np.isinf(sample["audio"]["array"]).any():
        print(f"Invalid audio data in sample: {sample['path']}")
        sample["audio"]["array"] = np.zeros_like(sample["audio"]["array"])

    return sample

train_dataset = [resample_audio(sample) for sample in train_dataset]
test_dataset = [resample_audio(sample) for sample in test_dataset]

import IPython.display as ipd
import numpy as np
import random

rand_int = random.randint(0, len(train_dataset))

print("Target text:", train_dataset[rand_int]['english_transcription'])
print("Input array shape:", np.asarray(train_dataset[rand_int]["audio"]["array"]).shape)
print("Sampling rate:", train_dataset[rand_int]["audio"]["sampling_rate"])
ipd.Audio(data=np.asarray(train_dataset[rand_int]["audio"]["array"]), autoplay=True, rate=16000)

from transformers import WhisperProcessor, WhisperForConditionalGeneration

processor = WhisperProcessor.from_pretrained("openai/whisper-base")
model = WhisperForConditionalGeneration.from_pretrained("openai/whisper-base")

def prepare_dataset(example):
    audio = example["audio"]
    input_features = processor.feature_extractor(
        audio["array"], sampling_rate=16000, return_tensors="pt"
    ).input_features[0]
    labels = processor.tokenizer(example["transcription"], return_tensors="pt").input_ids[0]
    return {"input_features": input_features, "labels": labels}

train_dataset = train_dataset.map(prepare_dataset, remove_columns=train_dataset.column_names)
eval_dataset = test_dataset.map(prepare_dataset, remove_columns=test_dataset.column_names)

from dataclasses import dataclass
from typing import Optional

@dataclass
class DataCollatorSpeechSeq2SeqWithPadding:
    processor: WhisperProcessor
    padding: bool = True
    max_length: Optional[int] = None
    pad_to_multiple_of: Optional[int] = None
    return_tensors: str = "pt"

    def __call__(self, features):
        input_features = [{"input_features": f["input_features"]} for f in features]
        label_features = [{"input_ids": f["labels"]} for f in features]

        batch = self.processor.feature_extractor.pad(
            input_features,
            padding=self.padding,
            max_length=self.max_length,
            pad_to_multiple_of=self.pad_to_multiple_of,
            return_tensors=self.return_tensors
        )

        labels_batch = self.processor.tokenizer.pad(
            label_features,
            padding=self.padding,
            max_length=self.max_length,
            pad_to_multiple_of=self.pad_to_multiple_of,
            return_tensors=self.return_tensors
        )

        labels = labels_batch["input_ids"].masked_fill(
            labels_batch["attention_mask"].ne(1), -100
        )

        batch["labels"] = labels
        return batch

data_collator = DataCollatorSpeechSeq2SeqWithPadding(processor=processor)

from transformers import TrainingArguments

training_args = TrainingArguments(
    output_dir="./whisper-finetuned",
    per_device_train_batch_size=4,
    gradient_accumulation_steps=4,
    learning_rate=1e-5,
    num_train_epochs=10,
    eval_strategy="steps",
    save_steps=500,
    eval_steps=500,
    logging_steps=100,
    fp16=True,
    save_total_limit=2,
    report_to="none",
)

from transformers import Trainer

trainer = Trainer(
    model=model,
    args=training_args,
    train_dataset=train_dataset,
    eval_dataset=eval_dataset,
    data_collator=data_collator,
    tokenizer=processor.tokenizer,
)

trainer.train()
trainer.evaluate()