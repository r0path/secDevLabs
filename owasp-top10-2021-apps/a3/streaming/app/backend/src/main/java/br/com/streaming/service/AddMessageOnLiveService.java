package br.com.streaming.service;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import br.com.streaming.dto.MessageDTO;
import br.com.streaming.model.Live;
import br.com.streaming.model.Message;
import br.com.streaming.model.User;
import br.com.streaming.repository.LiveRepository;
import br.com.streaming.repository.MessageRepository;
import br.com.streaming.repository.UserRepository;

@Service
public class AddMessageOnLiveService {
	
	@Autowired
	private LiveRepository liveRepository;

	@Autowired
	private UserRepository userRepository;
	
	@Autowired
	private MessageRepository messageRepository;

	@Transactional
	public Message add(Long id, MessageDTO messageDto) {
		// Ensure the target Live exists before creating or persisting related entities.
		Live live = liveRepository.findById(id).orElse(null);
		if (live == null) {
			// Do not persist messages (or users) when the live does not exist — avoids orphaned records.
			return null;
		}

		User user = userRepository.findByUsername(messageDto.username).orElse(null);
		if (user == null) {
			user = new User(messageDto.username, messageDto.username);
			userRepository.save(user);
		}

		Message message = new Message(user, messageDto.content);

		// Associate in-memory first, then persist. Transactional ensures both saves roll back on failure.
		live.addMessage(message);
		messageRepository.save(message);
		liveRepository.save(live);
		return message;
	}
}
